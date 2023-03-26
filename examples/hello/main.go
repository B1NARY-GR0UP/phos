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
	"log"

	"github.com/B1NARY-GR0UP/phos"
)

func main() {
	ph := phos.New[int](0)
	plusOne := func(ctx context.Context, data int) (int, error) {
		return data + 1, nil
	}
	ph.Handlers = append(ph.Handlers, plusOne)
	ph.In <- 0
	res := <-ph.Out
	log.Printf("res: %d", res.Data)
}
