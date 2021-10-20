module knative.dev/kn-plugin-event

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v1.1.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/thediveo/enumflag v0.10.0
	github.com/wavesoftware/go-ensure v1.0.0
	go.uber.org/zap v1.19.1
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	k8s.io/klog/v2 v2.20.0 // indirect
	knative.dev/client v0.26.1-0.20211019150534-534d91319f7d
	knative.dev/eventing v0.26.1-0.20211020051652-2db2eef82875
	knative.dev/hack v0.0.0-20211019034732-ced8ce706528
	knative.dev/networking v0.0.0-20211019132235-c8d647402afa
	knative.dev/pkg v0.0.0-20211019132235-ba2b2b1bf268
	knative.dev/serving v0.26.1-0.20211020110753-e4b3034b0df6
	sigs.k8s.io/yaml v1.3.0
)
