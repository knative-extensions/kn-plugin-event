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

package k8s

import "errors"

var (
	// ErrInvalidReference if given reference is invalid.
	ErrInvalidReference = errors.New("reference is invalid")

	// ErrNotFound if given reference do not point to any resource.
	ErrNotFound = errors.New("resource not found")

	// ErrNotAddressable if found resource isn't addressable.
	ErrNotAddressable = errors.New("resource isn't addressable")

	// ErrMoreThenOneFound if more then one resource has been found.
	ErrMoreThenOneFound = errors.New("more then one resource has been found")

	// ErrUnexcpected if something unexpected actually has happened.
	ErrUnexcpected = errors.New("something unexpected actually has happened")
)
