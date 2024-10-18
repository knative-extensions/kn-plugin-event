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
