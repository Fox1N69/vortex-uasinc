package service

import (
	"context"
	"fmt"
	"test-task/infra/k8s"
	"test-task/internal/models"
	"test-task/internal/repository"
	"test-task/pkg/util/logger"
	"time"
)

type ClientService interface {
	Create(client *models.Client) (int64, error)
	ClientByID(id int64) (*models.Client, error)
	Update(id int64, updateParams map[string]interface{}) error
	Delete(id int64) error
	Clients(ctx context.Context) ([]models.Client, error)
	AlgorithmStatuses() ([]models.AlgorithmStatus, error)
	UpdateAlgorithmStatus(id int64, status map[string]interface{}) error
	StartAlgorithmSync()
}

type clientService struct {
	repository  repository.ClientRepository
	k8sDeployer k8s.KubernetesDeployer
	log         logger.Logger
}

func NewClientService(clientRepo repository.ClientRepository, k8sDeployer k8s.KubernetesDeployer) ClientService {
	logger := logger.GetLogger()
	return &clientService{
		repository:  clientRepo,
		k8sDeployer: k8sDeployer,
		log:         logger,
	}
}

func (cs *clientService) Create(client *models.Client) (int64, error) {
	var algorithm models.AlgorithmStatus
	return cs.repository.Create(client, &algorithm)
}

func (cs *clientService) ClientByID(id int64) (*models.Client, error) {
	return cs.repository.ClientByID(id)
}

func (cs *clientService) Update(id int64, updateParams map[string]interface{}) error {
	return cs.repository.Update(id, updateParams)
}

func (cs *clientService) Delete(id int64) error {
	return cs.repository.Delete(id)
}

func (cs *clientService) Clients(ctx context.Context) ([]models.Client, error) {
	return cs.repository.Clients(ctx)
}

func (cs *clientService) AlgorithmStatuses() ([]models.AlgorithmStatus, error) {
	return cs.repository.AlgorithmStatuses()
}

func (cs *clientService) UpdateAlgorithmStatus(id int64, status map[string]interface{}) error {
	return cs.repository.UpdateAlgorithmStatus(id, status)
}

// StartAlgorithmSync initiates the algorithm synchronization process.
// This function starts a goroutine that synchronizes algorithms every 5 minute.
// A Ticker is used to trigger the synchronization at the specified intervals.
// When the function completes, the Ticker is stopped to release resources.
func (cs *clientService) StartAlgorithmSync() {
	const op = "service.client.StartAlgorithmSync"

	cs.log.Infof("%s: Starting synchronization process...", op)
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			cs.syncAlgorithms()
		}
	}()

	cs.log.Infof("%s: Synchronization process started", op)
}

// syncAlgorithms fetches clients from the database and synchronizes pods for each client based on their algorithm status.
func (cs *clientService) syncAlgorithms() {
	const op = "service.client.syncAlgorithms"

	clients, err := cs.repository.Clients(context.Background())
	if err != nil {
		cs.log.Errorf("%s: Failed to fetch clients from database: %v", op, err)
		return
	}

	for _, client := range clients {
		algoStatus, err := cs.repository.AlgorithmByClientID(context.Background(), client.ID)
		if err != nil {
			cs.log.Errorf("%s: Failed to fetch algorithm status for client %d: %v", op, client.ID, err)
			continue
		}
		cs.syncPodsForClient(client, *algoStatus)
	}
}

// syncPodsForClient synchronizes pods for a given client based on their algorithm status.
// It creates or deletes Kubernetes pods depending on the algorithm status flags VWAP, TWAP, and HFT.
// For each algorithm type, a pod is created if the corresponding flag is true in algoStatus;
// otherwise, the pod is deleted.
// Pod names are generated based on the client's ID and algorithm type (e.g., "vwap-123").
// If pod creation or deletion fails, an error is logged.
func (cs *clientService) syncPodsForClient(client models.Client, algoStatus models.AlgorithmStatus) {
	const op = "service.client.syncPodsForClient"

	// VWAP
	vwapPodName := fmt.Sprintf("vwap-%d", client.ID)
	if algoStatus.VWAP {
		if err := cs.k8sDeployer.CreatePod(vwapPodName, client.Image); err != nil {
			cs.log.Errorf("%s: Failed to deploy VWAP pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: VWAP pod deployed successfully for client %d", op, client.ID)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(vwapPodName); err != nil {
			cs.log.Errorf("%s: Failed to delete VWAP pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: VWAP pod deleted successfully for client %d", op, client.ID)
		}
	}

	// TWAP
	twapPodName := fmt.Sprintf("twap-%d", client.ID)
	if algoStatus.TWAP {
		if err := cs.k8sDeployer.CreatePod(twapPodName, client.Image); err != nil {
			cs.log.Errorf("%s: Failed to deploy TWAP pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: TWAP pod deployed successfully for client %d", op, client.ID)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(twapPodName); err != nil {
			cs.log.Errorf("%s: Failed to delete TWAP pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: TWAP pod deleted successfully for client %d", op, client.ID)
		}
	}

	// HFT
	hftPodName := fmt.Sprintf("hft-%d", client.ID)
	if algoStatus.HFT {
		if err := cs.k8sDeployer.CreatePod(hftPodName, client.Image); err != nil {
			cs.log.Errorf("%s: Failed to deploy HFT pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: HFT pod deployed successfully for client %d", op, client.ID)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(hftPodName); err != nil {
			cs.log.Errorf("%s: Failed to delete HFT pod for client %d: %v", op, client.ID, err)
		} else {
			cs.log.Infof("%s: HFT pod deleted successfully for client %d", op, client.ID)
		}
	}
}
