package configuration

import (
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/sender"
)

func memoizeKubeClients(delegate sender.CreateKubeClients) sender.CreateKubeClients {
	mem := kubeClientsMemoizer{delegate: delegate}
	return mem.computeClients
}

type kubeClientsMemoizer struct {
	delegate sender.CreateKubeClients
	result   k8s.Clients
}

func (m *kubeClientsMemoizer) computeClients(props *event.Properties) (k8s.Clients, error) {
	if m.result != nil {
		return m.result, nil
	}
	cl, err := m.delegate(props)
	if err != nil {
		return nil, err
	}
	m.result = cl
	return m.result, nil
}
