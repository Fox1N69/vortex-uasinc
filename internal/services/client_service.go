package service

import (
	"test-task/internal/models"
	"test-task/internal/repository"
	"time"
)

type ClientService interface {
	Create(client *models.Client) (int64, error)
	ClientByID(id int64) (*models.Client, error)
	Update(id int64, updateParams map[string]interface{}) error
	Delete(id int64) error
	Clients() ([]models.Client, error)
	PeriodicCheck()
	AlgorithmStatuses() ([]models.AlgorithmStatus, error)
	UpdateAlgorithmStatus(id int64, status map[string]interface{}) error
}

type clientService struct {
	repository repository.ClientRepository
}

func NewClientService(clientRepo repository.ClientRepository) ClientService {
	return &clientService{
		repository: clientRepo,
	}
}

func (cs *clientService) Create(client *models.Client) (int64, error) {
	return cs.repository.Create(client)
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

func (cs *clientService) Clients() ([]models.Client, error) {
	return cs.repository.Clients()
}

func (cs *clientService) AlgorithmStatuses() ([]models.AlgorithmStatus, error) {
	return cs.repository.AlgorithmStatuses()
}

func (cs *clientService) UpdateAlgorithmStatus(id int64, status map[string]interface{}) error {
	return cs.repository.UpdateAlgorithmStatus(id, status)
}

func (cs *clientService) PeriodicCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cs.checkAndDeployPods()
		}
	}
}

func (cs *clientService) checkAndDeployPods() {
	// Реализация проверки статусов алгоритмов и управления pod'ами
}
