/*
Copyright 2021 The Knative Authors
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

package tests

import (
	"fmt"
	"os"
)

// WithEnviron will execute a block of code with temporal environment set.
func WithEnviron(env map[string]string, body func() error) error {
	old := map[string]*string{}
	for k := range env {
		var vv *string
		if v, ok := os.LookupEnv(k); ok {
			vv = &v
		}
		old[k] = vv
	}
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("can't set env %v=%v: %w", k, v, err)
		}
	}
	defer func() {
		for k, v := range old {
			if v != nil {
				_ = os.Setenv(k, *v)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}()
	return body()
}
