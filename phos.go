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
	options *Options

	In  chan<- T
	Out <-chan Result[T]

	pool sync.Pool
	once sync.Once

	handlers []Handler[T]
	appendC  chan Handler[T]
	removeC  chan int
}

// Handler handles the data of PHOS channel
// TODO: remove context
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
		In:       in,
		Out:      out,
		handlers: make([]Handler[T], 0),
		options:  options,
		appendC:  make(chan Handler[T]),
		removeC:  make(chan int),
	}
	ph.pool.New = func() any {
		return make(chan Result[T])
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
	})
}

// Len return the number of handlers
// Note: This method is not concurrency safe
// A recommended way to use this method is to call it before calling Remove
// e.g. ph.Remove(ph.Len() - 1)
func (ph *Phos[T]) Len() int {
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

// Pause PHOS execution
func (ph *Phos[T]) Pause(ctx context.Context) {
	// TODO: implement me
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	ctx := ph.options.Ctx
	for {
		select {
		case handler, ok := <-ph.appendC:
			if !ok {
				return
			}
			ph.handlers = append(ph.handlers, handler)
		case index, ok := <-ph.removeC:
			if !ok {
				return
			}
			if index < 0 || index > len(ph.handlers)-1 {
				continue
			}
			copy(ph.handlers[index:], ph.handlers[index+1:])
			ph.handlers = ph.handlers[:len(ph.handlers)-1]
		case data, ok := <-in:
			if !ok {
				out <- ph.result(data, false, nil)
				return
			}
			receiver := ph.pool.Get().(chan Result[T])
			timer := time.NewTimer(ph.options.Timeout)
			go ph.doHandle(ctx, data, receiver)
			select {
			case <-timer.C:
				timer.Stop()
				if ph.options.ErrTimeoutFunc != nil {
					data = ph.options.ErrTimeoutFunc(ctx, data).(T)
				}
				out <- ph.result(data, true, timeoutError())
			case res := <-receiver:
				timer.Stop()
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
}

func (ph *Phos[T]) doHandle(ctx context.Context, data T, receiver chan Result[T]) {
	// TODO: 超时后没有正确释放 goroutine，是不是应该把超时逻辑放在这里
	var err error
	for _, handler := range ph.handlers {
		data, err = handler(ctx, data)
		if err != nil {
			if ph.options.ErrHandleFunc != nil {
				data = ph.options.ErrHandleFunc(ctx, data, err).(T)
			}
			receiver <- ph.result(data, true, handlerError(err))
			ph.pool.Put(receiver)
			return
		}
	}
	receiver <- ph.result(data, true, nil)
	ph.pool.Put(receiver)
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
