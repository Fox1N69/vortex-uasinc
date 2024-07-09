package manager

import (
	"test-task/infra"
)

type RepoManager interface {
}

type repoManager struct {
	infra infra.Infra
}

// NewRepoManager creates a new instance of RepoManager using the provided infrastructure.
func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}
