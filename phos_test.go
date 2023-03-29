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
	ph := New[int]()
	ph.Handlers = append(ph.Handlers, plusOne)
	ph.In <- 0
	res, ok := <-ph.Out
	assert.Equal(t, 1, res.Data)
	assert.True(t, ok)
	assert.Nil(t, res.Err)
}

// TODO: 即使一次性输入的数据超出缓冲区大小，也不会阻塞
// TODO: 如果取出的数据超出缓冲区大小，会阻塞而不是报错
func TestMultiHandlers(t *testing.T) {
	ph := New[int]()
	ph.Handlers = append(ph.Handlers, plusOne, plusOne, plusOne)
	ph.In <- 1 // 1 + 1 + 1 + 1 = 4
	ph.In <- 2 // 2 + 1 + 1 + 1 = 5
	ph.In <- 3 // 3 + 1 + 1 + 1 = 6
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, 4, res1.Data)
	assert.Equal(t, 5, res2.Data)
	assert.Equal(t, 6, res3.Data)
	assert.Nil(t, res1.Err)
	assert.Nil(t, res2.Err)
	assert.Nil(t, res3.Err)
}

func TestHandlersWithErr(t *testing.T) {
	ph := New[int]()
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	// Note:
	// The last handler will not be executed because of the error
	ph.In <- 1 // 1 + 111 + 1 = 113
	ph.In <- 2 // 2 + 111 + 1 = 114
	ph.In <- 3 // 3 + 111 + 1 = 115
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 113, res1.Data)
	assert.Equal(t, 114, res2.Data)
	assert.Equal(t, 115, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithZeroOption(t *testing.T) {
	// Note:
	// WithZero will make the result of the handler with error to be zero value of the type
	ph := New[int](WithZero())
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 0, res1.Data)
	assert.Equal(t, 0, res2.Data)
	assert.Equal(t, 0, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithErrHandleFuncOption(t *testing.T) {
	// Note:
	// WithErrHandleFunc will make the result of the handler with error to be the return value of the errHandleFunc
	ph := New[int](WithErrHandleFunc(plusSixSixSix))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1 // 1 + 1 + 111 + 666 = 779
	ph.In <- 2 // 2 + 1 + 111 + 666 = 780
	ph.In <- 3 // 3 + 1 + 111 + 666 = 781
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 779, res1.Data)
	assert.Equal(t, 780, res2.Data)
	assert.Equal(t, 781, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithZeroAndErrHandleFuncOption(t *testing.T) {
	// Note:
	// When WithZero and WithErrHandleFunc were enabled at the same time,
	// WithZero will overwrite the result of the WithErrHandleFunc option, that is,
	// the outputs will be the zero value of the appropriate type when error occur
	ph := New[int](WithZero(), WithErrHandleFunc(plusSixSixSix))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 0, res1.Data)
	assert.Equal(t, 0, res2.Data)
	assert.Equal(t, 0, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "plus one error", res1.Err.Error())
	assert.Equal(t, "plus one error", res2.Err.Error())
	assert.Equal(t, "plus one error", res3.Err.Error())
}

func TestHandlersWithTimeout(t *testing.T) {
	ph := New[int]()
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithSleep, plusOne)
	ph.In <- 10
	ph.In <- 20
	ph.In <- 30
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 10, res1.Data)
	assert.Equal(t, 20, res2.Data)
	assert.Equal(t, 30, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "phos error timeout", res1.Err.Error())
	assert.Equal(t, "phos error timeout", res2.Err.Error())
	assert.Equal(t, "phos error timeout", res3.Err.Error())
}

func TestHandlersWithTimeoutOption(t *testing.T) {
	ph := New[int](WithTimeout(time.Second * 5))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithSleep, plusOne)
	ph.In <- 10
	ph.In <- 20
	ph.In <- 30
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 10, res1.Data)
	assert.Equal(t, 20, res2.Data)
	assert.Equal(t, 30, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "phos error timeout", res1.Err.Error())
	assert.Equal(t, "phos error timeout", res2.Err.Error())
	assert.Equal(t, "phos error timeout", res3.Err.Error())
}

func TestHandlersWithErrTimeoutFuncOption(t *testing.T) {
	ph := New[int](WithErrTimeoutFunc(plusFiveFiveFive))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithSleep, plusOne)
	ph.In <- 10 // 10 + 555 = 565
	ph.In <- 20 // 20 + 555 = 575
	ph.In <- 30 // 30 + 555 = 585
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 565, res1.Data)
	assert.Equal(t, 575, res2.Data)
	assert.Equal(t, 585, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "phos error timeout", res1.Err.Error())
	assert.Equal(t, "phos error timeout", res2.Err.Error())
	assert.Equal(t, "phos error timeout", res3.Err.Error())
}

func TestHandlersWithCtxDoneFuncOption(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ph := New[int](WithContext(ctx), WithCtxDoneFunc(plusFiveFiveFive))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithSleep, plusOne)
	ph.In <- 10 // 10 + 555 = 565
	ph.In <- 20 // 20 + 555 = 575
	ph.In <- 30 // 30 + 555 = 585
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	assert.Equal(t, 565, res1.Data)
	assert.Equal(t, 575, res2.Data)
	assert.Equal(t, 585, res3.Data)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, "context deadline exceeded", res1.Err.Error())
	assert.Equal(t, "context deadline exceeded", res2.Err.Error())
	assert.Equal(t, "context deadline exceeded", res3.Err.Error())
}

func TestDefaultFuncOption(t *testing.T) {
	ph := New[int](WithDefaultFunc(logHelloPHOS))
	ph.Handlers = append(ph.Handlers, plusOne, plusOneWithErr, plusOne)
	time.Sleep(time.Microsecond * 100)
}

func plusOne(_ context.Context, data int) (int, error) {
	return data + 1, nil
}

func plusOneWithErr(_ context.Context, data int) (int, error) {
	return data + 111, errors.New("plus one error")
}

func plusSixSixSix(_ context.Context, data any, _ error) any {
	return data.(int) + 666
}

func plusFiveFiveFive(_ context.Context, data any) any {
	return data.(int) + 555
}

func plusOneWithSleep(_ context.Context, data int) (int, error) {
	time.Sleep(time.Second * 6)
	return data + 1, nil
}

func logHelloPHOS(_ context.Context) {
	log.Println("hello phos")
}
