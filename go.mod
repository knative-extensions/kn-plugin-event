module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/ghodss/yaml v1.0.0
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
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/client v0.29.1-0.20220217100713-af052088caa5
	knative.dev/eventing v0.29.1-0.20220217054812-7a48f4269b6f
	knative.dev/hack v0.0.0-20220218190734-a8ef7b67feec
	knative.dev/pkg v0.0.0-20220217155112-d48172451966
	knative.dev/reconciler-test v0.0.0-20220216192840-2c3291f210ce
	knative.dev/serving v0.29.1-0.20220217223834-b28062cdc4c7
	sigs.k8s.io/yaml v1.3.0
)
