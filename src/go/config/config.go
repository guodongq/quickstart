package config

import "github.com/spf13/viper"

type Config struct {
	GitHubApp *GitHubAppConfig
}

func NewConfigFromEnv() *Config {
	v := viper.New()

	return &Config{
		GitHubApp: &GitHubAppConfig{
			ClientID:   v.GetString(GitHubAppClientIDEnvKey),
			PrivateKey: []byte(v.GetString(GitHubAppPrivateKeyEnvKey)),
		},
	}
}
