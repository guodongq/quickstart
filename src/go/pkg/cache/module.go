package cache

import (
	"github.com/guodongq/quickstart/config"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"cache",
	fx.Provide(
		func(config *config.Config) *config.TTLConfig {
			return config.Cache.TTL
		},
		func(config *config.Config) *config.LRUConfig {
			return config.Cache.LRU
		},
	),
	fx.Provide(
		fx.Annotate(
			func(config *config.TTLConfig) Cache {
				return NewTTL(config.DefaultExpiration, config.EvictionInterval)
			},
			fx.As(new(Cache)),
			fx.ResultTags(`name:"ttl"`),
		),
		fx.Annotate(
			func(config *config.LRUConfig) Cache {
				return NewLRU(config.DefaultExpiration, config.EvictionInterval, config.MaxEntries)
			},
			fx.As(new(Cache)),
			fx.ResultTags(`name:"lru"`),
		),
	),
)
