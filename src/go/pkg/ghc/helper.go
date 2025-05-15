package ghc

import (
	"fmt"
	"github.com/google/go-github/v72/github"
	"github.com/guodongq/quickstart/config"
	"golang.org/x/oauth2"
	"net/http"
)

func createBaseClient(config *config.GitHubAppConfig) *github.Client {
	return github.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Source: NewCachedTokenSource(NewJWTGenerator(config.ClientID, config.PrivateKey)),
		},
	})
}

func createInstallationClient(token *github.InstallationToken) *github.Client {
	return github.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.ReuseTokenSource(nil, &installationTokenSource{token}),
		},
	})
}

type installationTokenSource struct{ token *github.InstallationToken }

func (s *installationTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: s.token.GetToken(),
		Expiry:      s.token.GetExpiresAt().Time,
		TokenType:   "Bearer",
	}, nil
}
func handleResponse(resp *github.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}
