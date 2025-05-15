package config

import (
	stderrors "errors"
	"github.com/guodongq/quickstart/pkg/errors"
)

type GitHubAppConfig struct {
	ClientID   string
	PrivateKey []byte
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
