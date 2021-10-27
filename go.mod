module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.11.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	github.com/wavesoftware/go-magetasks v0.5.2
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.21.4
	knative.dev/client v0.26.1-0.20211026094721-2939f36bbd88
	knative.dev/eventing v0.26.1-0.20211026055720-a20ff702ad22
	knative.dev/hack v0.0.0-20211026141922-a71c865b5f66
	knative.dev/networking v0.0.0-20211027012100-31eca9c93514
	knative.dev/pkg v0.0.0-20211026205101-a8fb29270197
	knative.dev/serving v0.26.1-0.20211026205200-8f02c277d0b3
	sigs.k8s.io/yaml v1.3.0
)

// FIXME: google/ko requires 0.22, remove when knative will work with 0.22+
replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
