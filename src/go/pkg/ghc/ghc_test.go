package ghc_test

import (
	"context"
	"github.com/guodongq/quickstart/config"
	"github.com/guodongq/quickstart/pkg/ghc"
	"testing"
)

const clientID = "Iv23li7PHtx4IztL7RyM"

const privateKey = `
`

const installationID int64 = 61632378

func TestGithubApps(t *testing.T) {
	var (
		client = ghc.NewGitHubAppClient(&config.GitHubAppConfig{
			ClientID:   clientID,
			PrivateKey: []byte(privateKey),
		})
		ctx = context.Background()
	)

	installations, err := client.ListInstallations(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("installations: %v", installations)

	installationID, err := client.GetInstallationID(ctx, "test-github-apps-0225")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("installationID: %d", installationID)
}
