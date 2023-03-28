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
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleHandler(t *testing.T) {
	ph := New[int](0)
	ph.Handlers = append(ph.Handlers, plusOne)
	ph.In <- 0
	res := <-ph.Out
	assert.Equal(t, 1, res.Data)
	assert.Nil(t, res.Err)
}

func TestMultiHandlers(t *testing.T) {
	ph := New[int](3)
	ph.Handlers = append(ph.Handlers, plusOne, plusOne, plusOne)
	ph.In <- 1 // 1 + 1 + 1 + 1 = 4
	ph.In <- 2 // 2 + 1 + 1 + 1 = 5
	ph.In <- 3 // 3 + 1 + 1 + 1 = 6
	res1 := <-ph.Out
	res2 := <-ph.Out
	res3 := <-ph.Out
	assert.Equal(t, 4, res1.Data)
	assert.Equal(t, 5, res2.Data)
	assert.Equal(t, 6, res3.Data)
	assert.Nil(t, res1.Err)
	assert.Nil(t, res2.Err)
	assert.Nil(t, res3.Err)
}

func TestHandlersWithErr(t *testing.T) {
	ph := New[int](3)
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	// Note:
	// The last handler will not be executed because of the error
	ph.In <- 1 // 1 + 111 + 1 = 113
	ph.In <- 2 // 2 + 111 + 1 = 114
	ph.In <- 3 // 3 + 111 + 1 = 115
	res1 := <-ph.Out
	res2 := <-ph.Out
	res3 := <-ph.Out
	assert.Equal(t, 113, res1.Data)
	assert.Equal(t, 114, res2.Data)
	assert.Equal(t, 115, res3.Data)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithZeroOption(t *testing.T) {
	// Note:
	// WithZero will make the result of the handler with error to be zero value of the type
	ph := New[int](3, WithZero())
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	res1 := <-ph.Out
	res2 := <-ph.Out
	res3 := <-ph.Out
	assert.Equal(t, 0, res1.Data)
	assert.Equal(t, 0, res2.Data)
	assert.Equal(t, 0, res3.Data)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithErrHandleFuncOption(t *testing.T) {
	// Note:
	// WithErrHandleFunc will make the result of the handler with error to be the return value of the errHandleFunc
	ph := New[int](3, WithErrHandleFunc(plusSixSixSix))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1 // 1 + 1 + 111 + 666 = 779
	ph.In <- 2 // 2 + 1 + 111 + 666 = 780
	ph.In <- 3 // 3 + 1 + 111 + 666 = 781
	res1 := <-ph.Out
	res2 := <-ph.Out
	res3 := <-ph.Out
	assert.Equal(t, 779, res1.Data)
	assert.Equal(t, 780, res2.Data)
	assert.Equal(t, 781, res3.Data)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithZeroAndErrHandleFuncOption(t *testing.T) {
	// Note:
	// When WithZero and WithErrHandleFunc were enabled at the same time,
	// WithZero will overwrite the result of the WithErrHandleFunc option, that is,
	// the outputs will be the zero value of the appropriate type when error occur
	ph := New[int](3, WithZero(), WithErrHandleFunc(plusSixSixSix))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	res1 := <-ph.Out
	res2 := <-ph.Out
	res3 := <-ph.Out
	assert.Equal(t, 0, res1.Data)
	assert.Equal(t, 0, res2.Data)
	assert.Equal(t, 0, res3.Data)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

// TODO: fix data race
func TestHandlersWithTimeout(t *testing.T) {
	ph := New[int](0)
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithSleep, plusOne)
	for i := 0; i < 50; i++ {
		ph.In <- 10
		res1 := <-ph.Out
		log.Println(res1)
		log.Println()
		ph.In <- 30
		res2 := <-ph.Out
		log.Println(res2)
		log.Println()
	}
}

func plusOne(_ context.Context, data int) (int, error) {
	return data + 1, nil
}

func plusOneWithErr(_ context.Context, data int) (int, error) {
	return data + 111, errors.New("plus one error")
}

func plusSixSixSix(_ context.Context, data any, err error) any {
	return data.(int) + 666
}

func plusOneWithSleep(_ context.Context, data int) (int, error) {
	time.Sleep(time.Second * 5)
	return data + 1, nil
}
