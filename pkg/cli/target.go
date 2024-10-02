package cli

import (
	"errors"
	"fmt"
	"net/url"

	"knative.dev/client/pkg/flags/sink"
	"knative.dev/kn-plugin-event/pkg/event"
)

var (
	// ErrUseToFlagIsRequired will be raised if user hasn't used --to flag.
	ErrUseToFlagIsRequired = errors.New("use --to flag is required")
	// ErrInvalidURLFormat will be raised if given URL is invalid.
	ErrInvalidURLFormat = errors.New("invalid URL format")
	// ErrInvalidToFormat will be raised if given addressable doesn't have a valid
	// expected format.
	ErrInvalidToFormat = errors.New("--to flag has invalid format")
)

// ValidateTarget will perform validation on App element of target.
func ValidateTarget(args *TargetArgs) error {
	if args.Sink == "" {
		return ErrUseToFlagIsRequired
	}
	ref, err := sink.Parse(args.Sink, "default", sink.ComputeWithDefaultMappings(nil))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToFormat, args.Sink)
	}
	if ref.Type() == sink.TypeReference {
		if ref.Name == "" {
			return fmt.Errorf("%w: %s", ErrInvalidToFormat, args.Sink)
		}
	}
	if ref.Type() == sink.TypeURL && !isValidAbsURL(args.Sink) {
		return fmt.Errorf("%w: %s", ErrInvalidURLFormat, args.Sink)
	}
	return validateAddressableURI(args.AddressableURI)
}

func validateAddressableURI(uri string) error {
	if len(uri) > 0 {
		_, err := url.ParseRequestURI(uri)
		if err != nil {
			return fmt.Errorf("--addressable-uri %s: %w: %w",
				uri, ErrInvalidURLFormat, err)
		}
	}
	return nil
}

func (a *App) createTarget(args TargetArgs, params *Params) (*event.Target, error) {
	mappings := sink.ComputeWithDefaultMappings(nil)
	if ref, err := sink.Parse(args.Sink, "default", mappings); err == nil && ref.Type() == sink.TypeURL {
		// a special case to avoid K8s connection if unnecessary
		return &event.Target{
			Reference:   ref,
			RelativeURI: args.AddressableURI,
		}, nil
	}
	var namespace string
	clients, kerr := a.Binding.NewKubeClients(params.Parse())
	if kerr != nil {
		return nil, kerr
	}
	namespace = clients.Namespace()

	ref, err := sink.Parse(args.Sink, namespace, mappings)
	if err != nil {
		return nil, err
	}
	return &event.Target{
		Reference:   ref,
		RelativeURI: args.AddressableURI,
	}, nil
}

func isValidAbsURL(uri string) bool {
	u, err := url.Parse(uri)
	return err == nil && u.Host != "" && u.IsAbs()
}
