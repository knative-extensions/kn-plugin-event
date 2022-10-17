package config

import (
	"context"

	"github.com/wavesoftware/go-magetasks/pkg/cache"
)

// Cache return a cache.Cache implementation based on context.Context.
func Cache() cache.Cache {
	return contextCache{}
}

type contextCache struct{}

func (c contextCache) Compute(
	key interface{},
	provider func() (interface{}, error),
) (interface{}, error) {
	value := fromContext(key)
	if value != nil {
		return value, nil
	}
	value, err := provider()
	if err != nil {
		return nil, err
	}
	saveInContext(key, value)
	return value, nil
}

func (c contextCache) Drop(key interface{}) interface{} {
	value := fromContext(key)
	saveInContext(key, nil)
	return value
}

func saveInContext(cacheKey interface{}, value interface{}) {
	WithContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cacheKey, value)
	})
}

func fromContext(cacheKey interface{}) interface{} {
	return Actual().Context.Value(cacheKey)
}
