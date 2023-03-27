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

import "errors"

var _ error = (*Error)(nil)

type Error struct {
	Err  error
	Type ErrorType
}

func (e *Error) Error() string {
	return e.Err.Error()
}

type ErrorType uint64

const (
	_ ErrorType = iota
	TimeoutErr
	HandleErr
	CtxErr
)

func newError(err error, t ErrorType) *Error {
	return &Error{
		Err:  err,
		Type: t,
	}
}

func timeoutError() *Error {
	return newError(errors.New("phos error timeout"), TimeoutErr)
}

func handleError(err error) *Error {
	return newError(err, HandleErr)
}

func ctxError(err error) *Error {
	return newError(err, CtxErr)
}
