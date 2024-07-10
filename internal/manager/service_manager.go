package manager

import (
	"sync"
	"test-task/infra"
	service "test-task/internal/services"
)

type ServiceManager interface {
	ClientService() service.ClientService
}

type serviceManager struct {
	infra infra.Infra
	repo  RepoManager
}

// NewServiceManager ...
func NewServiceManager(infra infra.Infra) ServiceManager {
	return &serviceManager{
		infra: infra,
		repo:  NewRepoManager(infra),
	}
}

var (
	clientServiceOnce sync.Once
	clientService     service.ClientService
)

func (sm *serviceManager) ClientService() service.ClientService {
	clientServiceOnce.Do(func() {
		clientRepo := sm.repo.ClientRepository()
		clientService = service.NewClientService(clientRepo)
	})

	return clientService
}
