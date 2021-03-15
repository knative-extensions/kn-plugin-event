package cli

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/cardil/kn-event/internal/event"
	clientutil "knative.dev/client/pkg/util"
	"knative.dev/pkg/apis"
)

var (
	// ErrCantUseBothToURLAndToFlags will be raised if user use both --to and
	// --to-url flags.
	ErrCantUseBothToURLAndToFlags = errors.New("can't use both --to and --to-url flags")
	// ErrUseToURLOrToFlagIsRequired will be raised if user didn't used --to or
	// --to-url flags.
	ErrUseToURLOrToFlagIsRequired = errors.New("use --to or --to-url flag is required")
	// ErrInvalidURLFormat will be raised if given URL is invalid.
	ErrInvalidURLFormat = errors.New("invalid URL format")
	// ErrInvalidToFormat will be raised if given addressable doesn't have valid
	// expected format.
	ErrInvalidToFormat = errors.New("--to flag needs to be in format " +
		"kind:apiVersion:name for named resources or " +
		"kind:apiVersion:labelKey1=value1,labelKey2=value2 for matching via " +
		"a label selector")
)

// ValidateTarget will perform validation on App element of target.
func ValidateTarget(args *TargetArgs) error {
	if args.URL == "" && args.Addressable == "" {
		return ErrUseToURLOrToFlagIsRequired
	}
	if args.URL != "" && args.Addressable != "" {
		return ErrCantUseBothToURLAndToFlags
	}
	if args.URL != "" {
		_, err := url.ParseRequestURI(args.URL)
		if err != nil {
			return fmt.Errorf("--to-url %w: %s", ErrInvalidURLFormat, err.Error())
		}
	}
	if args.Addressable != "" {
		// ref: https://regex101.com/r/TcxsLO/3
		r := regexp.MustCompile("([a-zA-Z0-9]+):([a-zA-Z0-9/.]+):([a-zA-Z0-9=,_-]+)")
		if !r.MatchString(args.Addressable) {
			return ErrInvalidToFormat
		}
	}
	_, err := url.ParseRequestURI(args.AddressableURI)
	if err != nil {
		return fmt.Errorf("--addressable-uri %w: %s", ErrInvalidURLFormat, err.Error())
	}
	return nil
}

func createTarget(args *TargetArgs, props *event.Properties) (*event.Target, error) {
	if args.Addressable != "" {
		ref, err := clientutil.ToTrackerReference(args.Addressable, args.Namespace)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidToFormat, err.Error())
		}
		uri, err := apis.ParseURL(args.AddressableURI)
		if err != nil {
			return nil, fmt.Errorf("--addressable-uri %w: %s", ErrInvalidURLFormat, err.Error())
		}
		return &event.Target{
			Type: event.TargetTypeAddressable,
			// FIXME: .Reference.Namespace and .SenderNamespace needs to be filled in if they eql ""
			AddressableVal: &event.AddressableSpec{
				Reference:       ref,
				URI:             uri,
				SenderNamespace: args.SenderNamespace,
			},
			Properties: props,
		}, nil
	}
	if args.URL != "" {
		u, err := url.Parse(args.URL)
		if err != nil {
			return nil, fmt.Errorf("--to-url %w: %s", ErrInvalidURLFormat, err.Error())
		}
		return &event.Target{
			Type:       event.TargetTypeReachable,
			URLVal:     u,
			Properties: props,
		}, nil
	}
	return nil, ErrUseToURLOrToFlagIsRequired
}
