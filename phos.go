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
)

// Phos short for Phosphophyllite
// PHOS is a channel with internal handler chain
type Phos[T any] struct {
	In  chan<- T
	Out <-chan Result[T]

	options  *Options
	handlers chan Handler[T]
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
		In:       in,
		Out:      out,
		options:  options,
		handlers: make(chan Handler[T]),
	}

	go ph.handle(in, out)

	return ph
}

func (ph *Phos[T]) AddHandler(h Handler[T]) {
	go func() {
		ph.handlers <- h
	}()
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	var handlers []Handler[T]
	for {
		select {
		case handler := <-ph.handlers:
			handlers = append(handlers, handler)

		case data, ok := <-in:
			if !ok {
				out <- Result[T]{
					Data: data,
					OK:   false,
					Err:  nil,
				}
				continue
			}

			notifier := make(chan Result[T])
			ctx, cancel := context.WithTimeout(ph.options.Ctx, ph.options.Timeout)

			go ph.executeHandlers(ctx, handlers, data, notifier)

			select {
			case result := <-notifier:
				cancel()
				out <- result
			case <-ctx.Done():
				cancel()
				out <- Result[T]{
					Data: data,
					OK:   ok,
					Err:  timeoutError(),
				}
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
			notifier <- Result[T]{
				Data: data,
				OK:   true,
				Err:  handleError(err),
			}
			return
		}
	}

	notifier <- Result[T]{
		Data: data,
		OK:   true,
		Err:  nil,
	}
}
