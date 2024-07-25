package e2e

import "knative.dev/kn-plugin-event/test/images"

func init() { //nolint:gochecknoinits
	// TODO: Remove the deprecatedResolver after CI is updated to use the new
	//       naming convention.
	deprecatedResolver := &images.EnvironmentalBasedResolver{
		Prefix: "test-images",
	}
	resolver := &images.EnvironmentalBasedResolver{
		Prefix: "client-plugin-event-test",
	}
	images.Resolvers = append(images.Resolvers, resolver, deprecatedResolver)
}
