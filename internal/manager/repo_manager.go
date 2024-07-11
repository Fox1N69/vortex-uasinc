package manager

import (
	"sync"
	"test-task/infra"
	"test-task/internal/repository"
)

type RepoManager interface {
	ClientRepository() repository.ClientRepository
}

type repoManager struct {
	infra infra.Infra
}

// NewRepoManager creates a new instance of RepoManager using the provided infrastructure.
func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	clientRepositoryOnce sync.Once
	clientRepository     repository.ClientRepository
)

func (rm *repoManager) ClientRepository() repository.ClientRepository {
	clientRepositoryOnce.Do(func() {
		clientRepository = repository.NewClientRepository(rm.infra.PSQLClient().DB)
	})
	return clientRepository
}
