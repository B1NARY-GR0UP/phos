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
	In          chan<- T
	Out         <-chan Result[T]
	handlerChan chan Handler[T]
	pool        sync.Pool
	options     *Options
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
		In:          in,
		Out:         out,
		handlerChan: make(chan Handler[T]),
		options:     options,
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
	close(ph.In)
	close(ph.handlerChan)
}

// Append add handler for PHOS to execute
func (ph *Phos[T]) Append(handlers ...Handler[T]) {
	for _, handler := range handlers {
		ph.handlerChan <- handler
	}
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	ctx := ph.options.Ctx
	handlers := make([]Handler[T], 0)
	for {
		select {
		case handler, ok := <-ph.handlerChan:
			if !ok {
				return
			}
			handlers = append(handlers, handler)
		case data, ok := <-in:
			if !ok {
				out <- ph.result(data, false, nil)
				continue
			}
			receiver := ph.pool.Get().(chan Result[T])
			timer := time.NewTimer(ph.options.Timeout)
			go ph.executeHandlers(ctx, handlers, data, receiver)
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

func (ph *Phos[T]) executeHandlers(ctx context.Context, handlers []Handler[T], data T, notifier chan Result[T]) {
	var err error
	for _, handler := range handlers {
		data, err = handler(ctx, data)
		if err != nil {
			if ph.options.ErrHandleFunc != nil {
				data = ph.options.ErrHandleFunc(ctx, data, err).(T)
			}
			notifier <- ph.result(data, true, handleError(err))
			ph.pool.Put(notifier)
			return
		}
	}
	notifier <- ph.result(data, true, nil)
	ph.pool.Put(notifier)
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
