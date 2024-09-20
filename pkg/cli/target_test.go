package cli_test

import (
	"errors"
	"testing"

	"knative.dev/kn-plugin-event/pkg/cli"
)

func TestValidateTarget(t *testing.T) {
	tests := []struct {
		name    string
		args    cli.TargetArgs
		wantErr error
	}{{
		name:    "empty is invalid",
		wantErr: cli.ErrUseToFlagIsRequired,
	}, {
		name: "valid URL",
		args: cli.TargetArgs{
			Sink:           "https://example.org",
			AddressableURI: "/",
		},
		wantErr: nil,
	}, {
		name: "invalid URL",
		args: cli.TargetArgs{
			Sink: "https://",
		},
		wantErr: cli.ErrInvalidURLFormat,
	}, {
		name: "invalid addressable URI",
		args: cli.TargetArgs{
			Sink:           "https://example.org",
			AddressableURI: "This is not an URI",
		},
		wantErr: cli.ErrInvalidURLFormat,
	}, {
		name: "valid addressable",
		args: cli.TargetArgs{
			Sink:           "service:serving.knative.dev/v1:showcase",
			AddressableURI: "/",
		},
		wantErr: nil,
	}, {
		name: "invalid sink",
		args: cli.TargetArgs{
			Sink:           "service::showcase",
			AddressableURI: "/",
		},
		wantErr: cli.ErrInvalidToFormat,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.ValidateTarget(&tt.args); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateTarget():\n   error = %#v\n wantErr = %#v", err, tt.wantErr)
			}
		})
	}
}
