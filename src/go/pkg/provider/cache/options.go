package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type CacheOptions struct {
	Ctx    context.Context
	Config bigcache.Config
}

func getDefaultCacheOptions() CacheOptions {
	return CacheOptions{
		Ctx:    context.Background(),
		Config: bigcache.DefaultConfig(10 * time.Minute),
	}
}

func WithCacheOptionsCtx(ctx context.Context) func(options *CacheOptions) {
	return func(options *CacheOptions) {
		options.Ctx = ctx
	}
}

func WithCacheOptionsConfig(config bigcache.Config) func(options *CacheOptions) {
	return func(options *CacheOptions) {
		options.Config = config
	}
}
