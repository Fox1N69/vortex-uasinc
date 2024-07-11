package service_test

import (
	"context"
	"test-task/internal/models"
	service "test-task/internal/services"
	"testing"
	"time"

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

type MockKubernetesDeployer struct {
	mock.Mock
}

func (m *MockKubernetesDeployer) CreatePod(name, image string) error {
	args := m.Called(name, image)
	return args.Error(0)
}

func (m *MockKubernetesDeployer) DeletePod(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockKubernetesDeployer) GetPodList() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func TestClientService_Create(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	client := &models.Client{ID: 1, ClientName: "Test Client"}
	mockRepo.On("Create", client).Return(int64(1), nil)

	id, err := service.Create(client)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestClientService_ClientByID(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	client := &models.Client{ID: 1, ClientName: "Test Client"}
	mockRepo.On("ClientByID", int64(1)).Return(client, nil)

	res, err := service.ClientByID(int64(1))

	assert.NoError(t, err)
	assert.Equal(t, client, res)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Update(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	updateParams := map[string]interface{}{"ClientName": "Updated Client"}
	mockRepo.On("Update", int64(1), updateParams).Return(nil)

	err := service.Update(int64(1), updateParams)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Delete(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := service.Delete(int64(1))

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestClientService_Clients(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

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
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	algorithm := &models.AlgorithmStatus{ID: 1, ClientID: 1, VWAP: true}
	mockRepo.On("CreateAlgorithm", algorithm).Return(int64(1), nil)

	id, err := service.CreateAlgorithm(algorithm)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockRepo.AssertExpectations(t)
}

func TestClientService_AlgorithmStatuses(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

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
	mockK8sDeployer := new(MockKubernetesDeployer)
	service := service.NewClientService(mockRepo, mockK8sDeployer)

	updateParams := map[string]interface{}{"VWAP": true}
	mockRepo.On("UpdateAlgorithmStatus", int64(1), updateParams).Return(nil)

	err := service.UpdateAlgorithmStatus(int64(1), updateParams)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStartAlgorithmSync(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockK8sDeployer := new(MockKubernetesDeployer)

	service := service.NewClientService(mockRepo, mockK8sDeployer)

	clients := []models.Client{
		{ID: 1, ClientName: "Client1"},
		{ID: 2, ClientName: "Client2"},
	}
	mockRepo.On("Clients", mock.Anything).Return(clients, nil)

	mockRepo.On("AlgorithmByClientID", mock.Anything, int64(1)).Return(&models.AlgorithmStatus{VWAP: true}, nil)
	mockRepo.On("AlgorithmByClientID", mock.Anything, int64(2)).Return(&models.AlgorithmStatus{VWAP: false}, nil)

	mockK8sDeployer.On("CreatePod", mock.Anything, mock.Anything).Return(nil)
	mockK8sDeployer.On("DeletePod", mock.Anything).Return(nil)

	go service.StartAlgorithmSync()

  time.Sleep(5 * time.Second)	

	mockRepo.AssertExpectations(t)
}
