package cache

// NoopCache do not do any caching. Probably, should be only used for testing
// purposes.
type NoopCache struct{}

func (n NoopCache) Compute(_ interface{}, provider func() (interface{}, error)) (interface{}, error) {
	return provider()
}

func (n NoopCache) Drop(_ interface{}) interface{} {
	return nil
}
