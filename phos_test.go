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
	"fmt"
	"testing"
)

func handlePlusOne(ctx context.Context, data int) (int, error) {
	return data + 1, nil
}

func handlePlusOneWithErr(ctx context.Context, data int) (int, error) {
	if data == 1 {
		return data + 111, errors.New("plus one error")
	}
	return data + 1, nil
}

func handleMulti3(ctx context.Context, data int) (int, error) {
	return data * 3, nil
}

func TestFeasibility(t *testing.T) {
	ph := New[int](1)
	ph.Handlers = append(ph.Handlers, handleMulti3)
	ph.In <- 13
	res, ok := <-ph.Out
	fmt.Println(res, ok)
}

func TestMultiHandlers(t *testing.T) {
	ph := New[int](3)
	ph.Handlers = append(ph.Handlers, handlePlusOne, handlePlusOne, handlePlusOne)
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	fmt.Println(<-ph.Out)
	fmt.Println(<-ph.Out)
	fmt.Println(<-ph.Out)
}

func TestSingleHandlerErr(t *testing.T) {
	ph := New[int](3)
	ph.Handlers = append(ph.Handlers, handlePlusOneWithErr)
	ph.In <- 1
	fmt.Println(<-ph.Out)
}

func TestSingleHandlerErrWithZero(t *testing.T) {
	ph := New[int](3, WithZero())
	ph.Handlers = append(ph.Handlers, handlePlusOneWithErr)
	ph.In <- 1
	ph.In <- 2
	fmt.Println(<-ph.Out)
	fmt.Println(<-ph.Out)
}

func TestMultiHandlerErr(t *testing.T) {
	ph := New[int](3)
	ph.Handlers = append(ph.Handlers, handlePlusOneWithErr, handlePlusOne, handlePlusOne)
	ph.In <- 1 // 1 + 111
	ph.In <- 2 // 2 + 1 + 1 + 1
	ph.In <- 3 // 3 + 1 + 1 + 1
	fmt.Println(<-ph.Out)
	fmt.Println(<-ph.Out)
	fmt.Println(<-ph.Out)
}
