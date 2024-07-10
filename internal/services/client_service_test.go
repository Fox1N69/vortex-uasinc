package service_test

/*
import (
	"test-task/internal/models"
	service "test-task/internal/services"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) AddClient(client models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

// Implement other methods similarly

func TestAddClient(t *testing.T) {
	mockClientRepo := new(MockClientRepository)
	mockAlgoRepo := new(MockAlgorithmStatusRepository)
	mockDeployer := new(MockDeployer)
	log := logrus.NewEntry(logrus.StandardLogger())

	cs := service.NewClientService(mockClientRepo, mockAlgoRepo, mockDeployer, log)

	client := models.Client{
		// Initialize client fields
	}

	mockClientRepo.On("AddClient", client).Return(nil)

	_, err := cs.Create(&client)
	assert.NoError(t, err)

	mockClientRepo.AssertExpectations(t)
}
*/
