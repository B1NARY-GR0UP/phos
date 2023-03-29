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
	In       chan<- T
	Out      <-chan Result[T]
	Handlers []Handler[T]
	pool     sync.Pool
	options  *Options
}

// Handler handles the data of PHOS channel
type Handler[T any] func(ctx context.Context, data T) (T, error)

// Result PHOS output result
type Result[T any] struct {
	Data T
	OK   bool
	Err  *Error
}

// New PHOS channel
func New[T any](opts ...Option) *Phos[T] {
	options := newOptions(opts...)
	in := make(chan T, 1)
	out := make(chan Result[T], 1)
	ph := &Phos[T]{
		In:      in,
		Out:     out,
		options: options,
	}
	ph.pool.New = func() any {
		return make(chan struct{})
	}
	go ph.handle(in, out)
	return ph
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	ctx := ph.options.Ctx
NEXT:
	select {
	case data := <-in:
		notifier := ph.pool.Get().(chan struct{})
		timer := time.NewTimer(ph.options.Timeout)
		go ph.executeHandlers(ctx, data, out, notifier)
		select {
		case <-timer.C:
			timer.Stop()
			if ph.options.ErrTimeoutFunc != nil {
				data = ph.options.ErrTimeoutFunc(ctx, data).(T)
			}
			ph.launch(out, data, timeoutError())
			goto NEXT
		case <-notifier:
			timer.Stop()
			goto NEXT
		case <-ctx.Done():
			timer.Stop()
			if ph.options.CtxDoneFunc != nil {
				data = ph.options.CtxDoneFunc(ctx, data).(T)
			}
			ph.launch(out, data, ctxError(ctx.Err()))
			goto NEXT
		}
	}
}

func (ph *Phos[T]) executeHandlers(ctx context.Context, data T, out chan Result[T], notifier chan struct{}) {
	var err error
	for _, handler := range ph.Handlers {
		data, err = handler(ctx, data)
		if err != nil {
			if ph.options.ErrHandleFunc != nil {
				data = ph.options.ErrHandleFunc(ctx, data, err).(T)
			}
			notifier <- struct{}{}
			ph.pool.Put(notifier)
			ph.launch(out, data, handleError(err))
			return
		}
	}
	notifier <- struct{}{}
	ph.pool.Put(notifier)
	ph.launch(out, data, nil)
}

func (ph *Phos[T]) launch(out chan Result[T], data T, err *Error) {
	if ph.options.Zero && err != nil {
		var zero T
		out <- Result[T]{
			Data: zero,
			OK:   true,
			Err:  err,
		}
	} else {
		out <- Result[T]{
			Data: data,
			OK:   true,
			Err:  err,
		}
	}
}
