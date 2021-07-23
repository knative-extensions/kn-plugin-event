package sender

import "knative.dev/kn-plugin-event/pkg/event"

// DefaultNamespace returns a default namespace of connected K8s cluster or
// error if such namespace can't be determined.
func (b *Binding) DefaultNamespace(props *event.Properties) (string, error) {
	clients, err := b.CreateKubeClients(props)
	if err != nil {
		return "", err
	}
	return clients.Namespace(), nil
}
