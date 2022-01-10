module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.11.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/spf13/cobra v1.2.1
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	github.com/wavesoftware/go-magetasks v0.6.0
	go.uber.org/zap v1.19.1
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/client v0.28.1-0.20220110123359-c6997da944bd
	knative.dev/eventing v0.28.1-0.20220107145225-eb4c06c8009d
	knative.dev/hack v0.0.0-20211222071919-abd085fc43de
	knative.dev/networking v0.0.0-20220107020122-0dbedcd88acf
	knative.dev/pkg v0.0.0-20220105211333-96f18522d78d
	knative.dev/serving v0.28.1-0.20220107170125-03091748d279
	sigs.k8s.io/yaml v1.3.0
)
