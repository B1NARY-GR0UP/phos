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

// defaultOptions default options for PHOS
var defaultOptions = Options{
	Ctx:            context.Background(),
	Zero:           false,
	Timeout:        time.Second * 3,
	ErrHandleFunc:  nil,
	ErrTimeoutFunc: nil,
	DefaultFunc:    nil,
}

// Option for PHOS
type Option func(o *Options)

// Options for PHOS
type Options struct {
	Ctx            context.Context
	Zero           bool
	Timeout        time.Duration
	ErrHandleFunc  ErrHandlerFunc
	ErrTimeoutFunc ErrTimeoutFunc
	DefaultFunc    DefaultFunc
}

type (
	ErrHandlerFunc func(ctx context.Context, data any, err error)
	ErrTimeoutFunc func(ctx context.Context, data any)
	DefaultFunc    func(ctx context.Context)
)

// NewOptions new options for PHOS
func NewOptions(opts ...Option) *Options {
	options := &Options{
		Ctx:            defaultOptions.Ctx,
		Zero:           defaultOptions.Zero,
		Timeout:        defaultOptions.Timeout,
		ErrHandleFunc:  defaultOptions.ErrHandleFunc,
		ErrTimeoutFunc: defaultOptions.ErrTimeoutFunc,
		DefaultFunc:    defaultOptions.DefaultFunc,
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

// WithZero will return zero value when handle error happened
func WithZero() Option {
	return func(o *Options) {
		o.Zero = true
	}
}

// WithTimeout will set timeout for handlers execution
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// WithErrHandleFunc will set error handle function which will be called when error happened
func WithErrHandleFunc(fn ErrHandlerFunc) Option {
	return func(o *Options) {
		o.ErrHandleFunc = fn
	}
}

// WithErrTimeoutFunc will set error timeout function which will be called when timeout happened
func WithErrTimeoutFunc(fn ErrTimeoutFunc) Option {
	return func(o *Options) {
		o.ErrTimeoutFunc = fn
	}
}

// WithDefaultFunc will set default function which will be called when no data in channel
func WithDefaultFunc(fn DefaultFunc) Option {
	return func(o *Options) {
		o.DefaultFunc = fn
	}
}
