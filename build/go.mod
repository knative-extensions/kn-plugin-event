module knative.dev/kn-plugin-event/build

go 1.16

require (
	github.com/magefile/mage v1.11.0
	github.com/wavesoftware/go-magetasks v0.5.2
	knative.dev/kn-plugin-event v0.0.0
)

replace knative.dev/kn-plugin-event => ../

// FIXME: google/ko requires 0.22, remove when knative will work with 0.22+
replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.4
