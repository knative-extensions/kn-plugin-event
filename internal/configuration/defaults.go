package configuration

import (
	"knative.dev/kn-plugin-event/internal/event"
	"knative.dev/kn-plugin-event/internal/k8s"
	"knative.dev/kn-plugin-event/internal/sender"
)

func senderBinding() sender.Binding {
	return sender.Binding{
		CreateKubeClients:     k8s.CreateKubeClient,
		CreateJobRunner:       k8s.CreateJobRunner,
		CreateAddressResolver: k8s.CreateAddressResolver,
	}
}

func eventsBinding(binding sender.Binding) event.Binding {
	return event.Binding{
		CreateSender: binding.New,
	}
}
