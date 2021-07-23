package configuration

import (
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
)

func senderBinding() sender.Binding {
	return sender.Binding{
		CreateKubeClients:     memoizeKubeClients(k8s.CreateKubeClient),
		CreateJobRunner:       k8s.CreateJobRunner,
		CreateAddressResolver: k8s.CreateAddressResolver,
	}
}

func eventsBinding(binding sender.Binding) event.Binding {
	return event.Binding{
		CreateSender:     binding.New,
		DefaultNamespace: binding.DefaultNamespace,
	}
}
