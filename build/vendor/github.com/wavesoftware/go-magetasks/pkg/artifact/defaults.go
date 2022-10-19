package artifact

import "github.com/wavesoftware/go-magetasks/config"

// ConfigureDefaults will configure default builders and publishers to be used.
func ConfigureDefaults() {
	config.DefaultBuilders = append(config.DefaultBuilders,
		BinaryBuilder{},
		KoBuilder{},
	)
	config.DefaultPublishers = append(config.DefaultPublishers,
		ListPublisher{},
		KoPublisher{},
	)
}
