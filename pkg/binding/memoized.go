package binding

import (
	"knative.dev/kn-plugin-event/pkg/k8s"
)

func memoizeKubeClients(delegate k8s.NewKubeClients) k8s.NewKubeClients {
	mem := kubeClientsMemoizer{delegate: delegate}
	return mem.computeClients
}

type kubeClientsMemoizer struct {
	delegate k8s.NewKubeClients
	result   k8s.Clients
}

func (m *kubeClientsMemoizer) computeClients(params *k8s.Configurator) (k8s.Clients, error) {
	if m.result != nil {
		return m.result, nil
	}
	cl, err := m.delegate(params)
	if err != nil {
		return nil, err
	}
	m.result = cl
	return m.result, nil
}
