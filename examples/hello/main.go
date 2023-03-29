// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/phos"
)

func plusOne(_ context.Context, data int) (int, error) {
	return data + 1, nil
}

func main() {
	ph := phos.New[int]()
	ph.Handlers = append(ph.Handlers, plusOne)
	// TODO: bug1: close(ph.In) 后，ph.Out 会一直取出零值数据，并且 ok 为 true
	// TODO: bug2: 不传入数据，ph.Out 却会取出经过 handler 的数据
	//close(ph.In)
	res4, ok4 := <-ph.Out
	fmt.Println(res4, ok4)
}
