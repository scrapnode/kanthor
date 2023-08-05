package cache

import (
	"context"
	"errors"
	"time"
)

func Warp[T any](cache Cache, ctx context.Context, key string, ttl time.Duration, handler func() (*T, error)) (*T, error) {
	entry, err := cache.Get(ctx, key)
	if err == nil {
		return Unmarshal[T](entry)
	}

	// we only accept ErrEntryNotFound
	if !errors.Is(err, ErrEntryNotFound) {
		return nil, err
	}

	// cache miss
	data, err := handler()
	if err != nil {
		return nil, err
	}
	entry, err = Marshal(data)
	if err != nil {
		return nil, err
	}

	// return the data with error from cache setting
	return data, cache.Set(ctx, key, entry, ttl)
}
