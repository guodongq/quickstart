package config

import (
	stderrors "errors"
	"github.com/guodongq/quickstart/pkg/errors"
	"github.com/spf13/viper"
)

type (
	Config struct {
		GitHubApp *GitHubAppConfig
		Academy   *Academy
	}

	Academy struct {
		BaseURL string
	}

	GitHubAppConfig struct {
		ClientID   string
		PrivateKey []byte
	}
)

func NewConfigFromEnv() *Config {
	v := viper.New()
	v.SetDefault(AcademyBaseURLEnvKey, "/api/academy")

	return &Config{
		GitHubApp: &GitHubAppConfig{
			ClientID:   v.GetString(GitHubAppClientIDEnvKey),
			PrivateKey: []byte(v.GetString(GitHubAppPrivateKeyEnvKey)),
		},
		Academy: &Academy{
			BaseURL: v.GetString(AcademyBaseURLEnvKey),
		},
	}
}

func (c *GitHubAppConfig) Validate() error {
	if len(c.ClientID) == 0 {
		return errors.BadRequestError(stderrors.New("client id is required"))
	}

	if len(c.PrivateKey) == 0 {
		return errors.BadRequestError(stderrors.New("private key is required"))
	}

	return nil
}
