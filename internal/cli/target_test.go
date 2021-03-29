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

package cli_test

import (
	"errors"
	"testing"

	"knative.dev/kn-plugin-event/internal/cli"
)

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name    string
		args    *cli.TargetArgs
		wantErr error
	}{{
		name:    "empty is invalid",
		args:    &cli.TargetArgs{},
		wantErr: cli.ErrUseToURLOrToFlagIsRequired,
	}, {
		name: "valid URL",
		args: &cli.TargetArgs{
			URL:            "http://example.org",
			AddressableURI: "/",
		},
		wantErr: nil,
	}, {
		name: "invalid URL",
		args: &cli.TargetArgs{
			URL: "foo.html",
		},
		wantErr: cli.ErrInvalidURLFormat,
	}, {
		name: "invalid addressable URI",
		args: &cli.TargetArgs{
			URL:            "http://example.org",
			AddressableURI: "This is not an URI",
		},
		wantErr: cli.ErrInvalidURLFormat,
	}, {
		name: "valid addressable",
		args: &cli.TargetArgs{
			Addressable:    "service:serving.knative.dev/v1:showcase",
			AddressableURI: "/",
		},
		wantErr: nil,
	}, {
		name: "invalid addressable",
		args: &cli.TargetArgs{
			Addressable:    "service::showcase",
			AddressableURI: "/",
		},
		wantErr: cli.ErrInvalidToFormat,
	}, {
		name: "both URL and addressable aren't valid",
		args: &cli.TargetArgs{
			URL:         "https://example.org/",
			Addressable: "service:serving.knative.dev/v1:showcase",
		},
		wantErr: cli.ErrCantUseBothToURLAndToFlags,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.ValidateTarget(tt.args); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateTarget():\n   error = %#v\n wantErr = %#v", err, tt.wantErr)
			}
		})
	}
}
