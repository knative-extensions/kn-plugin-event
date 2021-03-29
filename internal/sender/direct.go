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

package sender

import (
	"context"
	"fmt"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/wavesoftware/go-ensure"
)

type directSender struct {
	url url.URL
}

func (d *directSender) Send(ce cloudevents.Event) error {
	c, err := cloudevents.NewDefaultClient()
	ensure.NoError(err)

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.TODO(), d.url.String())

	// Send that Event.
	err = c.Send(ctx, ce)
	if !cloudevents.IsACK(err) {
		return fmt.Errorf("%v: %w", ErrCouldntBeSent, err)
	}

	return nil
}
