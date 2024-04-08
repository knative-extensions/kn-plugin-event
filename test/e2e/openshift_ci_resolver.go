package e2e

import "knative.dev/kn-plugin-event/test/images"

func init() { //nolint:gochecknoinits
	images.Resolvers = append(images.Resolvers, &images.EnvironmentalBasedResolver{
		Prefix: "test-images",
	})
}
