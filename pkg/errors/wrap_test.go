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

package errors_test

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/errors"
)

var (
	errFoo = errors.New("foo")
	errBar = errors.New("bar")
)

func TestWrap(t *testing.T) {
	err := errors.Wrap(errFoo, errBar)
	assert.ErrorIs(t, err, errFoo)
	assert.ErrorIs(t, err, errBar)
	assert.Error(t, err, "bar: foo")
}

func TestCause(t *testing.T) {
	ferr := fmt.Errorf("%w: fizz", errFoo)
	err := errors.Wrap(ferr, errBar)
	got := errors.Cause(err)
	assert.Equal(t, got, ferr)
	assert.Equal(t, errors.Cause(ferr), nil)
}

func TestUnwrapAll(t *testing.T) {
	ferr := fmt.Errorf("%w: fizz", errFoo)
	err := errors.Wrap(ferr, errBar)
	got := errors.UnwrapAll(err)
	assert.Assert(t, len(got) == 2)
	assert.Equal(t, got[0], errBar)
	assert.ErrorIs(t, got[1], errFoo)
}
