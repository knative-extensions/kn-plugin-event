module knative.dev/kn-plugin-event

go 1.16

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.11.0
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/prometheus/common v0.30.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	github.com/wavesoftware/go-magetasks v0.5.2
	go.uber.org/zap v1.19.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.21.4 // indirect
	k8s.io/client-go v0.21.4
	knative.dev/client v0.25.1
	knative.dev/eventing v0.25.2
	knative.dev/hack v0.0.0-20210622141627-e28525d8d260
	knative.dev/networking v0.0.0-20210903132258-9d8ab8618e5f
	knative.dev/pkg v0.0.0-20210902173607-844a6bc45596
	knative.dev/serving v0.25.1
	sigs.k8s.io/yaml v1.3.0
)

// FIXME: google/ko requires 0.22, remove when knative will work with 0.22+
replace k8s.io/apimachinery v0.22.2 => k8s.io/apimachinery v0.21.4
