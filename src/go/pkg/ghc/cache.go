package ghc

import (
	"github.com/google/go-github/v72/github"
	"sync"
)

type TokenCache interface {
	Get(installationID int64) (*github.InstallationToken, bool)
	Set(installationID int64, token *github.InstallationToken)
}

type MemoryTokenCache struct {
	tks map[int64]*github.InstallationToken
	mu  sync.RWMutex
}

func NewMemoryTokenCache() *MemoryTokenCache {
	return &MemoryTokenCache{
		tks: make(map[int64]*github.InstallationToken),
	}
}

func (c *MemoryTokenCache) Get(installationID int64) (*github.InstallationToken, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tk, exists := c.tks[installationID]
	return tk, exists
}

func (c *MemoryTokenCache) Set(installationID int64, token *github.InstallationToken) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tks[installationID] = token
}
