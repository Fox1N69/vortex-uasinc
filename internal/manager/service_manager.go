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

// ClientService returns an instance of the client service.
// It lazily initializes the service on the first call using the client repository and Kubernetes deployer from the infrastructure.
func (sm *serviceManager) ClientService() service.ClientService {
	clientServiceOnce.Do(func() {
		clientRepo := sm.repo.ClientRepository()
		clientService = service.NewClientService(clientRepo, sm.infra.KubernetesDeployer())
	})

	return clientService
}
