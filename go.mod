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
	github.com/wavesoftware/go-magetasks v0.5.1
	go.uber.org/zap v1.19.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.21.4
	knative.dev/client v0.25.1-0.20210921123237-09f14b105638
	knative.dev/eventing v0.25.1-0.20210920134735-f031ba23b23d
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/networking v0.0.0-20210914225408-69ad45454096
	knative.dev/pkg v0.0.0-20210919202233-5ae482141474
	knative.dev/serving v0.25.1-0.20210920201536-4a26f1daa58a
	sigs.k8s.io/yaml v1.3.0
)

// FIXME: google/ko requires 0.22, remove when knative will work with 0.22+
replace k8s.io/apimachinery v0.22.2 => k8s.io/apimachinery v0.21.4
