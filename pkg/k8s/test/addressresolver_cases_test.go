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

package test_test

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/k8s/test"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestResolveAddressTestCases(t *testing.T) {
	tcs := test.ResolveAddressTestCases("default")
	assert.Check(t, len(tcs) >= 1)
}

func TestEnsureResolveAddress(t *testing.T) {
	tcs := test.ResolveAddressTestCases("default")
	tc := tcs[0]
	ctx := context.TODO()
	doNothing := func(testing.TB) {}
	test.EnsureResolveAddress(ctx, t, tc, func() (k8s.Clients, func(tb testing.TB)) {
		return &tests.FakeClients{TB: t, Objects: tc.Objects}, doNothing
	})
}
