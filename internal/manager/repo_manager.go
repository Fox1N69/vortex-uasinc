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

// ClientRepository returns an instance of the client repository.
// It lazily initializes the repository on the first call using the PSQLClient from the infrastructure.
func (rm *repoManager) ClientRepository() repository.ClientRepository {
	clientRepositoryOnce.Do(func() {
		clientRepository = repository.NewClientRepository(rm.infra.PSQLClient().DB)
	})
	return clientRepository
}
