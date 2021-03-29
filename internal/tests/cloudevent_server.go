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
	"context"
	"net/http/httptest"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// WithCloudEventsServer is a testing utility that help by starting a CloudEvents
// HTTP server which can catch a sent event.
func WithCloudEventsServer(test func(serverURL url.URL) error) (*cloudevents.Event, error) {
	var ce *cloudevents.Event
	receive := func(ctx context.Context, event cloudevents.Event) {
		ce = &event
	}
	ctx := context.Background()
	protocol, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, err
	}
	handler, err := cloudevents.NewHTTPReceiveHandler(ctx, protocol, receive)
	if err != nil {
		return nil, err
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	u, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}
	err = test(*u)
	return ce, err
}
