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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	// TimeoutError
	timeoutErr := timeoutError()
	assert.Equal(t, TimeoutErr, timeoutErr.Type)
	assert.Equal(t, "phos error timeout", timeoutErr.Err.Error())
	// HandleError
	handleErr := handleError(errors.New("handle error"))
	assert.Equal(t, HandleErr, handleErr.Type)
	assert.Equal(t, "handle error", handleErr.Err.Error())
	// CtxError
	ctxErr := ctxError(errors.New("ctx error"))
	assert.Equal(t, CtxErr, ctxErr.Type)
	assert.Equal(t, "ctx error", ctxErr.Err.Error())
}
