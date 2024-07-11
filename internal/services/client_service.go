package service

import (
	"context"
	"fmt"
	"test-task/infra/k8s"
	"test-task/internal/models"
	"test-task/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
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
}

func NewClientService(clientRepo repository.ClientRepository, k8sDeployer k8s.KubernetesDeployer) ClientService {
	return &clientService{
		repository:  clientRepo,
		k8sDeployer: k8sDeployer,
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
	log.Infof("Starting synchronization process...")
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			cs.syncAlgorithms()
		}
	}()
}

func (cs *clientService) syncAlgorithms() {
	clients, err := cs.repository.Clients(context.Background())
	if err != nil {
		log.Errorf("Failed to fetch clients from database: %v", err)
		return
	}

	for _, client := range clients {
		algoStatus, err := cs.repository.AlgorithmByClientID(context.Background(), client.ID)
		if err != nil {
			log.Errorf("Failed to fetch algorithm status for client %d: %v", client.ID, err)
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
	// VWAP
	vwapPodName := fmt.Sprintf("vwap-%d", client.ID)
	if algoStatus.VWAP {
		if err := cs.k8sDeployer.CreatePod(vwapPodName, client.Image); err != nil {
			log.Errorf("Failed to deploy VWAP pod for client %d: %v", client.ID, err)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(vwapPodName); err != nil {
			log.Errorf("Failed to delete VWAP pod for client %d: %v", client.ID, err)
		}
	}

	// TWAP
	twapPodName := fmt.Sprintf("twap-%d", client.ID)
	if algoStatus.TWAP {
		if err := cs.k8sDeployer.CreatePod(twapPodName, client.Image); err != nil {
			log.Errorf("Failed to deploy TWAP pod for client %d: %v", client.ID, err)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(twapPodName); err != nil {
			log.Errorf("Failed to delete TWAP pod for client %d: %v", client.ID, err)
		}
	}

	// HFT
	hftPodName := fmt.Sprintf("hft-%d", client.ID)
	if algoStatus.HFT {
		if err := cs.k8sDeployer.CreatePod(hftPodName, client.Image); err != nil {
			log.Errorf("Failed to deploy HFT pod for client %d: %v", client.ID, err)
		}
	} else {
		if err := cs.k8sDeployer.DeletePod(hftPodName); err != nil {
			log.Errorf("Failed to delete HFT pod for client %d: %v", client.ID, err)
		}
	}
}
