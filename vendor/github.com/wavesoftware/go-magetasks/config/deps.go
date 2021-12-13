package config

type Dependencies interface {
	Configurator
	Installs() []string
	merge(other Dependencies)
}

func NewDependencies(deps ...string) Dependencies {
	s := make(map[string]bool, len(deps))
	for _, dep := range deps {
		s[dep] = exists
	}
	return dependencies{
		set: s,
	}
}

const exists = true

type dependencies struct {
	set map[string]bool
}

func (d dependencies) Installs() []string {
	keys := make([]string, len(d.set))

	i := 0
	for k := range d.set {
		keys[i] = k
		i++
	}

	return keys
}

func (d dependencies) Configure(cfg Configurable) {
	cfg.Config().Dependencies.merge(d)
}

func (d dependencies) merge(other Dependencies) {
	for _, dep := range other.Installs() {
		d.set[dep] = exists
	}
}
