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
	"fmt"
	"github.com/B1NARY-GR0UP/phos"
)

func main() {
	ph := phos.New[int]()
	ph.In <- 1
	ph.In <- 2
	ph.In <- 3
	res1, ok1 := <-ph.Out
	res2, ok2 := <-ph.Out
	res3, ok3 := <-ph.Out
	fmt.Println(res1, ok1)
	fmt.Println(res2, ok2)
	fmt.Println(res3, ok3)
	// TODO: 不会出现 false 的情况
	res4, ok4 := <-ph.Out
	fmt.Println(res4, ok4)
	res5, ok5 := <-ph.Out
	fmt.Println(res5, ok5)
	res6, ok6 := <-ph.Out
	fmt.Println(res6, ok6)
}
