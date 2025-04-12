package cache

import (
	"github.com/allegro/bigcache/v3"
	"github.com/guodongq/quickstart/pkg/util/provider"
)

type Cache struct {
	provider.AbstractProvider

	options CacheOptions
	cache   *bigcache.BigCache
}

func New(optionFuncs ...func(*CacheOptions)) *Cache {
	defaultOptions := getDefaultCacheOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &Cache{
		options: *options,
	}
}

func (p *Cache) Init() error {
	cache, err := bigcache.New(p.options.Ctx, p.options.Config)
	if err != nil {
		return err
	}

	p.cache = cache
	return nil
}

func (p *Cache) Close() error {
	if p.cache == nil {
		return p.AbstractProvider.Close()
	}
	return p.cache.Close()
}

func (p *Cache) Get(key string) ([]byte, error) {
	return p.cache.Get(key)
}

func (p *Cache) Set(key string, value []byte) error {
	return p.cache.Set(key, value)
}
