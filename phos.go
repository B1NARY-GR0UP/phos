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
	// In read data from users
	In chan<- T
	// Out write data to users
	Out <-chan Result[T]
	// Handlers is the handler chain
	Handlers []Handler[T]
	// options for PHOS
	options *Options
}

// Handler handles the data of channel
type Handler[T any] func(ctx context.Context, data T) (T, error)

// Result PHOS output result
type Result[T any] struct {
	Data    T
	Err     error
	Timeout bool
}

// New PHOS channel
// TODO: support async handle?
// TODO: support channel status?
func New[T any](cap int, opts ...Option) *Phos[T] {
	options := NewOptions(opts...)
	in := make(chan T, cap)
	out := make(chan Result[T], cap)
	ph := &Phos[T]{
		In:      in,
		Out:     out,
		options: options,
	}
	go ph.handle(in, out)
	return ph
}

func (ph *Phos[T]) handle(in chan T, out chan Result[T]) {
	ctx := ph.options.Ctx
	for {
	NEXT:
		select {
		case data := <-in:
			var err error
			// TODO: handle timeout
			for _, handler := range ph.Handlers {
				data, err = handler(ctx, data)
				if err != nil {
					if ph.options.ErrHandleFunc != nil {
						ph.options.ErrHandleFunc(ctx, data, err)
					}
					if ph.options.Zero {
						// TODO: consider the case of there is a timeout and an exception
						ph.zero(out, err, false)
					} else {
						ph.value(out, data, err, false)
					}
					goto NEXT
				}
			}
			ph.value(out, data, err, false)
		default:
			if ph.options.DefaultFunc != nil {
				ph.options.DefaultFunc(ctx)
			}
		}
	}
}

func (ph *Phos[T]) zero(out chan Result[T], err error, timeout bool) {
	var zero T
	out <- Result[T]{
		Data:    zero,
		Err:     err,
		Timeout: timeout,
	}
}

func (ph *Phos[T]) value(out chan Result[T], data T, err error, timeout bool) {
	out <- Result[T]{
		Data:    data,
		Err:     err,
		Timeout: timeout,
	}
}
