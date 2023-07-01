// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except In compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to In writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package phos

import (
	"context"
	"sync"
	"time"
)

// Phos short for Phosphophyllite
// PHOS is a channel with internal handler chain
type Phos[T any] struct {
	In  chan<- T
	Out <-chan Result[T]

	handlers []Handler[T]

	options *Options

	once sync.Once
	mu   sync.RWMutex
	wg   sync.WaitGroup

	appendC  chan Handler[T]
	removeC  chan int
	receiveC chan Result[T]
	closeC   chan struct{}
}

// Handler handles the data of PHOS channel
type Handler[T any] func(ctx context.Context, data T) (T, error)

// Result PHOS output result
type Result[T any] struct {
	Data T
	// Note: You should use the OK of Result rather than the second return value of PHOS Out channel
	OK  bool
	Err *Error
}

// New PHOS channel
func New[T any](opts ...Option) *Phos[T] {
	options := newOptions(opts...)
	in := make(chan T, 1)
	out := make(chan Result[T], 1)
	ph := &Phos[T]{
		options:  options,
		In:       in,
		Out:      out,
		handlers: make([]Handler[T], 0),
		appendC:  make(chan Handler[T]),
		removeC:  make(chan int),
		receiveC: make(chan Result[T]),
		closeC:   make(chan struct{}),
	}
	go ph.handle(in, out)
	return ph
}

// Close PHOS channel
// Note: You should not close In channel manually before or after calling Close
func (ph *Phos[T]) Close() {
	ph.once.Do(func() {
		close(ph.In)
		close(ph.appendC)
		close(ph.removeC)
		<-ph.closeC
		close(ph.receiveC)
	})
}

// Len return the number of handlers
func (ph *Phos[T]) Len() int {
	ph.mu.RLock()
	defer ph.mu.RUnlock()
	return len(ph.handlers)
}

// Append add handler for PHOS to execute
func (ph *Phos[T]) Append(handlers ...Handler[T]) {
	for _, handler := range handlers {
		ph.appendC <- handler
	}
}

// Remove remove handler from PHOS
func (ph *Phos[T]) Remove(index int) {
	ph.removeC <- index
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	defer close(ph.closeC)
	ctx := ph.options.Ctx
LOOP:
	for {
		select {
		case handler, ok := <-ph.appendC:
			if !ok {
				break LOOP
			}
			ph.handlers = append(ph.handlers, handler)
		case index, ok := <-ph.removeC:
			if !ok {
				break LOOP
			}
			if index < 0 || index > len(ph.handlers)-1 {
				continue
			}
			copy(ph.handlers[index:], ph.handlers[index+1:])
			ph.handlers = ph.handlers[:len(ph.handlers)-1]
		case data, ok := <-in:
			if !ok {
				out <- ph.result(data, false, nil)
				break LOOP
			}
			timer := time.NewTimer(ph.options.Timeout)
			go ph.doHandle(ctx, data)
			select {
			case <-timer.C:
				timer.Stop()
				if ph.options.ErrTimeoutFunc != nil {
					data = ph.options.ErrTimeoutFunc(ctx, data).(T)
				}
				out <- ph.result(data, true, timeoutError())
			case res, ok := <-ph.receiveC:
				timer.Stop()
				if !ok {
					break LOOP
				}
				out <- res
			case <-ctx.Done():
				timer.Stop()
				if ph.options.ErrDoneFunc != nil {
					data = ph.options.ErrDoneFunc(ctx, data, ctx.Err()).(T)
				}
				out <- ph.result(data, true, ctxError(ctx.Err()))
			}
		}
	}
	ph.wg.Wait()
}

func (ph *Phos[T]) doHandle(ctx context.Context, data T) {
	past := time.Now()
	launch := func(err *Error) {
		if time.Now().After(past.Add(ph.options.Timeout)) {
			return
		}
		select {
		case ph.receiveC <- ph.result(data, true, err):
		default:
		}
	}
	ph.wg.Add(1)
	defer ph.wg.Done()
	var err error
	for _, handler := range ph.handlers {
		data, err = handler(ctx, data)
		if err != nil {
			if ph.options.ErrHandleFunc != nil {
				data = ph.options.ErrHandleFunc(ctx, data, err).(T)
			}
			launch(handlerError(err))
			return
		}
	}
	launch(nil)
}

func (ph *Phos[T]) result(data T, ok bool, err *Error) Result[T] {
	if ph.options.Zero && err != nil {
		var zero T
		return Result[T]{
			Data: zero,
			OK:   ok,
			Err:  err,
		}
	} else {
		return Result[T]{
			Data: data,
			OK:   ok,
			Err:  err,
		}
	}
}
