module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/Azure/go-autorest/autorest v0.11.17 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.10 // indirect
	github.com/cloudevents/sdk-go/v2 v2.7.0
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
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.21.4
	knative.dev/client v0.28.1-0.20211222151818-2c9a236917da
	knative.dev/eventing v0.28.1-0.20211222204918-d8297456d455
	knative.dev/hack v0.0.0-20211222071919-abd085fc43de
	knative.dev/networking v0.0.0-20211223013028-62388a5f2853
	knative.dev/pkg v0.0.0-20211216142117-79271798f696
	knative.dev/serving v0.28.1-0.20211221064617-c69f92cdfce7
	sigs.k8s.io/yaml v1.3.0
)

// FIXME: google/ko requires 0.22, remove when knative will work with 0.22+
replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
