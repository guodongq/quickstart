package proxy

import "github.com/spf13/viper"

type EnvConfig struct {
	enabled   bool
	debug     bool
	port      int
	endpoint  string
	targetURL string
}

func LoadEnvConfig(envPrefix string) EnvConfig {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(envPrefix)

	v.SetDefault("ENABLED", true)
	v.SetDefault("DEBUG", false)
	v.SetDefault("PORT", 4040)
	v.SetDefault("ENDPOINT", "/")
	v.SetDefault("TARGET_URL", "http://localhost:8080")

	return EnvConfig{
		enabled:   v.GetBool("ENABLED"),
		debug:     v.GetBool("DEBUG"),
		port:      v.GetInt("PORT"),
		endpoint:  v.GetString("ENDPOINT"),
		targetURL: v.GetString("TARGET_URL"),
	}
}

func (e EnvConfig) Options() func(*ProxyOptions) {
	return func(options *ProxyOptions) {
		options.Enabled = e.enabled
		options.Debug = e.debug
		options.Port = e.port
		options.Endpoint = e.endpoint
		options.TargetURL = e.targetURL
	}
}
