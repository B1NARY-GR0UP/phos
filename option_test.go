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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	errHandleFunc := func(ctx context.Context, data any, err error) any {
		return nil
	}
	errTimeoutFunc := func(ctx context.Context, data any) any {
		return nil
	}
	errDoneFunc := func(ctx context.Context, data any, err error) any {
		return nil
	}
	options := newOptions(
		WithContext(context.TODO()),
		WithZero(),
		WithTimeout(time.Second*5),
		WithErrHandleFunc(errHandleFunc),
		WithErrTimeoutFunc(errTimeoutFunc),
		WithErrDoneFunc(errDoneFunc),
	)
	assert.Equal(t, context.TODO(), options.Ctx)
	assert.True(t, options.Zero)
	assert.Equal(t, time.Second*5, options.Timeout)
	assert.Equal(t, fmt.Sprintf("%p", errHandleFunc), fmt.Sprintf("%p", options.ErrHandleFunc))
	assert.Equal(t, fmt.Sprintf("%p", errTimeoutFunc), fmt.Sprintf("%p", options.ErrTimeoutFunc))
	assert.Equal(t, fmt.Sprintf("%p", errDoneFunc), fmt.Sprintf("%p", options.ErrDoneFunc))
}

func TestDefaultOptions(t *testing.T) {
	options := newOptions()
	assert.Equal(t, context.Background(), options.Ctx)
	assert.False(t, options.Zero)
	assert.Equal(t, time.Second*3, options.Timeout)
	assert.Nil(t, options.ErrHandleFunc)
	assert.Nil(t, options.ErrTimeoutFunc)
	assert.Nil(t, options.ErrDoneFunc)
}
