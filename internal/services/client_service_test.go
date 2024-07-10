package service_test

import (
	"context"
	"test-task/infra/k8s"
	"test-task/internal/models"
	service "test-task/internal/services"
	"test-task/pkg/util/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(client *models.Client) (int64, error) {
	args := m.Called(client)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockClientRepository) ClientByID(id int64) (*models.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Client), args.Error(1)
}

func (m *MockClientRepository) Update(id int64, updateParams map[string]interface{}) error {
	args := m.Called(id, updateParams)
	return args.Error(0)
}

func (m *MockClientRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClientRepository) Clients(ctx context.Context) ([]models.Client, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Client), args.Error(1)
}

func (m *MockClientRepository) CreateAlgorithm(algorithm *models.AlgorithmStatus) (int64, error) {
	args := m.Called(algorithm)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockClientRepository) AlgorithmStatuses() ([]models.AlgorithmStatus, error) {
	args := m.Called()
	return args.Get(0).([]models.AlgorithmStatus), args.Error(1)
}

func (m *MockClientRepository) UpdateAlgorithmStatus(id int64, status map[string]interface{}) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockClientRepository) AlgorithmByClientID(ctx context.Context, clientID int64) (*models.AlgorithmStatus, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(*models.AlgorithmStatus), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func TestClientService_Create(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	client := &models.Client{ID: 1, ClientName: "Test Client"}
	mockRepo.On("Create", client).Return(int64(1), nil)

	id, err := service.Create(client)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestClientService_ClientByID(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.Logger{}, mockK8sDeployer)

	client := &models.Client{ID: 1, ClientName: "Test Client"}
	mockRepo.On("ClientByID", int64(1)).Return(client, nil)

	res, err := service.ClientByID(int64(1))

	assert.NoError(t, err)
	assert.Equal(t, client, res)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Update(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	updateParams := map[string]interface{}{"ClientName": "Updated Client"}
	mockRepo.On("Update", int64(1), updateParams).Return(nil)

	err := service.Update(int64(1), updateParams)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Delete(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := service.Delete(int64(1))

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Clients(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	clients := []models.Client{
		{ID: 1, ClientName: "Test Client 1"},
		{ID: 2, ClientName: "Test Client 2"},
	}
	mockRepo.On("Clients", mock.Anything).Return(clients, nil)

	res, err := service.Clients(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, clients, res)
	mockRepo.AssertExpectations(t)
}

func TestClientService_CreateAlgorithm(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	algorithm := &models.AlgorithmStatus{ID: 1, ClientID: 1, VWAP: true}
	mockRepo.On("CreateAlgorithm", algorithm).Return(int64(1), nil)

	id, err := service.CreateAlgorithm(algorithm)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestClientService_AlgorithmStatuses(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	algorithms := []models.AlgorithmStatus{
		{ID: 1, ClientID: 1, VWAP: true},
		{ID: 2, ClientID: 2, VWAP: false},
	}
	mockRepo.On("AlgorithmStatuses").Return(algorithms, nil)

	res, err := service.AlgorithmStatuses()

	assert.NoError(t, err)
	assert.Equal(t, algorithms, res)
	mockRepo.AssertExpectations(t)
}

func TestClientService_UpdateAlgorithmStatus(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := k8s.NewKubernetesDeployer()
	service := service.NewClientService(mockRepo, logger.GetLogger(), mockK8sDeployer)

	updateParams := map[string]interface{}{"VWAP": true}
	mockRepo.On("UpdateAlgorithmStatus", int64(1), updateParams).Return(nil)

	err := service.UpdateAlgorithmStatus(int64(1), updateParams)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
