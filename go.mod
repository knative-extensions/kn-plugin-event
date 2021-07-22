module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/fatih/color v1.10.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.2.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/magefile/mage v1.10.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	github.com/wavesoftware/go-magetasks v0.4.3
	go.uber.org/zap v1.18.1
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	knative.dev/client v0.24.1-0.20210721192038-a7df7c76a0d1
	knative.dev/eventing v0.24.1-0.20210720191312-003a8a2d4ef5
	knative.dev/hack v0.0.0-20210622141627-e28525d8d260
	knative.dev/networking v0.0.0-20210721164339-e3c72d213a1a
	knative.dev/pkg v0.0.0-20210715175632-d9b7180af6f2
	knative.dev/serving v0.24.1-0.20210721140446-55a613957a4f
	sigs.k8s.io/structured-merge-diff/v4 v4.1.0 // indirect
)

// TODO: unpin for k8s 0.21+, see: https://github.com/knative/client/pull/1209
replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
