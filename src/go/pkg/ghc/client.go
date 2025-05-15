package ghc

import (
	"context"
	"fmt"
	"github.com/google/go-github/v72/github"
	"github.com/guodongq/quickstart/config"
	"sync"
	"time"
)

// GitHubAppClient defines the interface for GitHub App operations.
type GitHubAppClient interface {
	ListInstallations(ctx context.Context) ([]*github.Installation, error)
	GetInstallationID(ctx context.Context, installationOrg string) (int64, error)
	GetInstallationClient(ctx context.Context, installationID int64) (*github.Client, error)
	GetInstallationAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error)
	ListOrgRepositories(ctx context.Context, installationID int64, owner string) ([]*github.Repository, error)
	GetRepository(ctx context.Context, installationID int64, owner, repo string) (*github.Repository, error)
	GetDefaultBranch(ctx context.Context, installationID int64, owner, repo string) (string, error)
	GetBranches(ctx context.Context, installationID int64, owner, repo string) ([]*github.Branch, error)
	DoesFileExist(ctx context.Context, installationID int64, owner, repo, filePath string) (*github.RepositoryContent, bool, error)
	GetIssues(ctx context.Context, installationID int64, owner, repo string) ([]*github.Issue, error)
	CreateIssue(ctx context.Context, installationID int64, owner, repo string, issue *github.IssueRequest) (*github.Issue, error)
	CreateComment(ctx context.Context, installationID int64, owner string, repo string, issueNumber int, comment *github.IssueComment) (*github.IssueComment, error)
}

// ClientConfig holds configuration for the client.
type ClientConfig struct {
	TokenRefreshThreshold time.Duration
}

// clientImpl is the implementation of GitHubAppClient.
type clientImpl struct {
	config       *config.GitHubAppConfig
	baseClient   *github.Client
	tokenCache   TokenCache
	clientConfig ClientConfig
}

// NewGitHubAppClient creates a new GitHub App client.
func NewGitHubAppClient(
	config *config.GitHubAppConfig,
	opts ...func(*ClientConfig),
) GitHubAppClient {
	cfg := ClientConfig{
		TokenRefreshThreshold: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &clientImpl{
		config:       config,
		baseClient:   createBaseClient(config),
		tokenCache:   NewMemoryTokenCache(),
		clientConfig: cfg,
	}
}

func (c *clientImpl) ListInstallations(ctx context.Context) ([]*github.Installation, error) {
	mgr := installationManager{client: c.baseClient}
	return mgr.ListAll(ctx)
}

func (c *clientImpl) GetInstallationID(ctx context.Context, installationOrg string) (int64, error) {
	installations, err := c.ListInstallations(ctx)
	if err != nil {
		return 0, err
	}

	for _, installation := range installations {
		if *installation.AppSlug == installationOrg {
			return *installation.ID, nil
		}
	}

	err = fmt.Errorf("installation not found for org: %s", installationOrg)
	return 0, err
}

type installationManager struct {
	client *github.Client
}

func (m *installationManager) ListAll(ctx context.Context) ([]*github.Installation, error) {
	var installations []*github.Installation
	opts := &github.ListOptions{Page: 1, PerPage: 100}

	for {
		result, resp, err := m.client.Apps.ListInstallations(ctx, opts)
		if err := handleResponse(resp, err); err != nil {
			return nil, fmt.Errorf("list installations failed: %w", err)
		}

		installations = append(installations, result...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return installations, nil
}

func (m *installationManager) ValidateExists(ctx context.Context, installationID int64) error {
	_, resp, err := m.client.Apps.GetInstallation(ctx, installationID)
	return handleResponse(resp, err)
}

func (c *clientImpl) GetInstallationClient(ctx context.Context, installationID int64) (*github.Client, error) {
	token, err := c.GetInstallationAccessToken(ctx, installationID)
	if err != nil {
		return nil, err
	}
	return createInstallationClient(token), nil
}

func (c *clientImpl) GetInstallationAccessToken(ctx context.Context, installationID int64) (*github.InstallationToken, error) {
	mgr := tokenManager{
		client:     c.baseClient,
		tokenCache: c.tokenCache,
		cfg:        c.clientConfig,
	}
	if err := (&installationManager{client: c.baseClient}).ValidateExists(ctx, installationID); err != nil {
		return nil, fmt.Errorf("installation validation failed: %w", err)
	}

	return mgr.getInstallationToken(ctx, installationID)
}

type tokenManager struct {
	client     *github.Client
	tokenCache TokenCache
	refreshMux sync.Mutex
	cfg        ClientConfig
}

func (m *tokenManager) getInstallationToken(ctx context.Context, installationID int64) (*github.InstallationToken, error) {
	if token, exists := m.tokenCache.Get(installationID); exists && m.isTokenValid(token) {
		return token, nil
	}

	m.refreshMux.Lock()
	defer m.refreshMux.Unlock()

	// double check if token is still
	if token, exists := m.tokenCache.Get(installationID); exists && m.isTokenValid(token) {
		return token, nil
	}

	token, err := m.createInstallationToken(ctx, installationID)
	if err != nil {
		return nil, err
	}

	m.tokenCache.Set(installationID, token)
	return token, nil
}

func (m *tokenManager) isTokenValid(token *github.InstallationToken) bool {
	return time.Now().Add(m.cfg.TokenRefreshThreshold).Before(token.ExpiresAt.Time)
}

func (m *tokenManager) createInstallationToken(ctx context.Context, installationID int64) (*github.InstallationToken, error) {
	token, resp, err := m.client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err := handleResponse(resp, err); err != nil {
		return nil, fmt.Errorf("create token failed: %w", err)
	}
	return token, nil
}

func (c *clientImpl) ListOrgRepositories(ctx context.Context, installationID int64, orgName string) ([]*github.Repository, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var repos []*github.Repository
	opts := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{Page: 1}}

	for {
		result, resp, err := client.Repositories.ListByOrg(ctx, orgName, opts)
		if err := handleResponse(resp, err); err != nil {
			return nil, fmt.Errorf("list repositories failed: %w", err)
		}

		repos = append(repos, result...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return repos, nil
}

func (c *clientImpl) GetRepository(ctx context.Context, installationID int64, orgName, repoName string) (*github.Repository, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}

	result, resp, err := client.Repositories.Get(ctx, orgName, repoName)
	if err := handleResponse(resp, err); err != nil {
		return nil, fmt.Errorf("list repositories failed: %w", err)
	}

	return result, nil
}

func (c *clientImpl) GetDefaultBranch(ctx context.Context, installationID int64, orgName, repoName string) (string, error) {
	repository, err := c.GetRepository(ctx, installationID, orgName, repoName)
	if err != nil {
		return "", err
	}
	return repository.GetDefaultBranch(), nil
}

func (c *clientImpl) GetBranches(ctx context.Context, installationID int64, orgName, repoName string) ([]*github.Branch, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}
	var branches []*github.Branch
	opts := &github.BranchListOptions{ListOptions: github.ListOptions{Page: 1}}

	for {
		result, resp, err := client.Repositories.ListBranches(ctx, orgName, repoName, opts)
		if err := handleResponse(resp, err); err != nil {
			return nil, fmt.Errorf("list branches failed: %w", err)
		}

		branches = append(branches, result...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return branches, nil
}

func (c *clientImpl) DoesFileExist(ctx context.Context, installationID int64, orgName, repoName, filePath string) (*github.RepositoryContent, bool, error) {
	defaultBranch, err := c.GetDefaultBranch(ctx, installationID, orgName, repoName)
	if err != nil {
		return nil, false, err
	}

	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, false, err
	}

	fileContent, _, response, err := client.Repositories.GetContents(ctx, orgName, repoName, filePath, &github.RepositoryContentGetOptions{Ref: defaultBranch})
	if err := handleResponse(response, err); err != nil {
		return nil, false, fmt.Errorf("get file content failed: %w", err)
	}

	return fileContent, true, nil
}

func (c *clientImpl) GetIssues(ctx context.Context, installationID int64, orgName, repoName string) ([]*github.Issue, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}

	var issues []*github.Issue
	opts := &github.IssueListByRepoOptions{ListOptions: github.ListOptions{Page: 1}}

	for {
		result, resp, err := client.Issues.ListByRepo(ctx, orgName, repoName, opts)
		if err := handleResponse(resp, err); err != nil {
			return nil, fmt.Errorf("list issues failed: %w", err)
		}

		issues = append(issues, result...)
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return issues, nil
}

func (c *clientImpl) CreateIssue(ctx context.Context, installationID int64, owner, repo string, request *github.IssueRequest) (*github.Issue, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}

	issue, resp, err := client.Issues.Create(ctx, owner, repo, request)
	if err := handleResponse(resp, err); err != nil {
		return nil, fmt.Errorf("create issue failed: %w", err)
	}

	return issue, nil
}

func (c *clientImpl) CreateComment(ctx context.Context, installationID int64, owner string, repo string, issueNumber int, comment *github.IssueComment) (*github.IssueComment, error) {
	client, err := c.GetInstallationClient(ctx, installationID)
	if err != nil {
		return nil, err
	}

	comment, resp, err := client.Issues.CreateComment(ctx, owner, repo, issueNumber, comment)
	if err := handleResponse(resp, err); err != nil {
		return nil, fmt.Errorf("create comment failed: %w", err)
	}
	return comment, nil
}
