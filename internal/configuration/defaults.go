package configuration

import (
	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/k8s"
	"github.com/cardil/kn-event/internal/sender"
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
