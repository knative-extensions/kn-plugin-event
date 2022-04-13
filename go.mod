module knative.dev/kn-plugin-event

go 1.16

require (
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-containerregistry v0.8.1-0.20220219142810-1571d7fdc46e
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.11.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.3.0
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	github.com/wavesoftware/go-magetasks v0.6.0
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.1.0
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/client v0.30.2-0.20220413081508-bb7fd738212c
	knative.dev/eventing v0.30.1-0.20220413070808-1b90ef335226
	knative.dev/hack v0.0.0-20220411131823-6ffd8417de7c
	knative.dev/pkg v0.0.0-20220412134708-e325df66cb51
	knative.dev/reconciler-test v0.0.0-20220412165608-994f0c3fab62
	knative.dev/serving v0.30.1-0.20220413003907-2de1474b55ba
	sigs.k8s.io/yaml v1.3.0
)
