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

func hello(_ context.Context, data string) (string, error) {
	return data + " PHOS", nil
}

func main() {
	ph := phos.New[string]()
	ph.Handlers = append(ph.Handlers, hello)
	ph.In <- "BINARY"
	res := <-ph.Out
	fmt.Println(res.Data)
}
