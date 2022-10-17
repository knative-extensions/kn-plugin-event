package cache

// Cache allows computing values only once.
type Cache interface {
	// Compute will compute value only once using provider and save it in cache,
	// so subsequent calls will return value directly from cache. In case of
	// provider returning error, this error will be returned, but not saved in
	// cache.
	Compute(key interface{}, provider func() (interface{}, error)) (interface{}, error)
	// Drop will drop value from cache only if given key is found. Otherwise, it
	// will return nil.
	Drop(key interface{}) interface{}
}
