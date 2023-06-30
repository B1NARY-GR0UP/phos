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
	"time"
)

var defaultOptions = Options{
	Ctx:            context.Background(),
	Zero:           false,
	Timeout:        3 * time.Second,
	ErrHandleFunc:  nil,
	ErrTimeoutFunc: nil,
	ErrDoneFunc:    nil,
}

// Option for PHOS
type Option func(o *Options)

// Options for PHOS
type Options struct {
	Ctx            context.Context
	Zero           bool
	Timeout        time.Duration
	ErrHandleFunc  ErrHandleFunc
	ErrTimeoutFunc ErrTimeoutFunc
	ErrDoneFunc    ErrDoneFunc
}

// TODO: remove useless context
type (
	ErrHandleFunc  func(ctx context.Context, data any, err error) any
	ErrTimeoutFunc func(ctx context.Context, data any) any
	ErrDoneFunc    func(ctx context.Context, data any, err error) any
)

func newOptions(opts ...Option) *Options {
	options := &Options{
		Ctx:            defaultOptions.Ctx,
		Zero:           defaultOptions.Zero,
		Timeout:        defaultOptions.Timeout,
		ErrHandleFunc:  defaultOptions.ErrHandleFunc,
		ErrTimeoutFunc: defaultOptions.ErrTimeoutFunc,
		ErrDoneFunc:    defaultOptions.ErrDoneFunc,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithContext will set context for PHOS
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Ctx = ctx
	}
}

// WithZero will set zero value for return when error happened
func WithZero() Option {
	return func(o *Options) {
		o.Zero = true
	}
}

// WithTimeout will set timeout for the handler chain execution (not just for each handler)
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithErrHandleFunc will set error handle function for PHOS which will be called when handle error happened
func WithErrHandleFunc(fn ErrHandleFunc) Option {
	return func(o *Options) {
		o.ErrHandleFunc = fn
	}
}

// WithErrTimeoutFunc will set error timeout function for PHOS which will be called when timeout error happened
func WithErrTimeoutFunc(fn ErrTimeoutFunc) Option {
	return func(o *Options) {
		o.ErrTimeoutFunc = fn
	}
}

// WithErrDoneFunc will set err done function for PHOS which will be called when ctx done happened
// Note: You should use it will WithContext, otherwise it will not work
func WithErrDoneFunc(fn ErrDoneFunc) Option {
	return func(o *Options) {
		o.ErrDoneFunc = fn
	}
}
