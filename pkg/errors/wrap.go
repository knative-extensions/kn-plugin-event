/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package errors

import (
	"errors"
	"fmt"
)

// Wrap an error with provided wrap error. Will check if the error is already
// of given type to prevent over-wrapping.
func Wrap(err, wrap error) error {
	if errors.Is(err, wrap) {
		return err
	}
	return fmt.Errorf("%w: %w", wrap, err)
}

// UnwrapAll will get all the wrapped errors of a given one, regardless if a
// single error was wrapped or multiple.
func UnwrapAll(err error) []error {
	cause := errors.Unwrap(err)
	if cause != nil {
		return []error{cause}
	}
	eg, ok := err.(multipleErrors)
	if !ok {
		return nil
	}
	return append(([]error)(nil), eg.Unwrap()...)
}

// Cause will determine the cause error of the given one, by returning the
// second wrapped error.
func Cause(err error) error {
	errs := UnwrapAll(err)
	switch len(errs) {
	case 0, 1:
		return nil
	default:
		return errs[1]
	}
}

type multipleErrors interface {
	Unwrap() []error
}
