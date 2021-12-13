package cli

// PluginVersionOutput is a struct that is used to output project version in
// machine readable format.
type PluginVersionOutput struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Image   string `json:"image" yaml:"image"`
}
