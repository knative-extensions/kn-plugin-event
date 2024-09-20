package binding

import (
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
)

// CliApp creates the configured cli.App to work with.
func CliApp() *cli.App {
	return &cli.App{
		Binding: eventsBinding(senderBinding()),
	}
}

// IcsApp creates the configured ics.App to work with.
func IcsApp() *ics.App {
	return &ics.App{
		Binding: eventsBinding(senderBinding()),
	}
}

func senderBinding() sender.Binding {
	return sender.Binding{
		NewKubeClients:     memoizeKubeClients(k8s.NewClients),
		NewJobRunner:       k8s.NewJobRunner,
		NewAddressResolver: k8s.NewAddressResolver,
	}
}

func eventsBinding(binding sender.Binding) event.Binding {
	return event.Binding{
		CreateSender:   binding.New,
		NewKubeClients: binding.NewKubeClients,
	}
}
