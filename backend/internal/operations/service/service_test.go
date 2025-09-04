package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	leasingModels "github.com/rendley/vegshare/backend/internal/leasing/models"
	leasingRepository "github.com/rendley/vegshare/backend/internal/leasing/repository"
	operationsModels "github.com/rendley/vegshare/backend/internal/operations/models"
	operationsRepository "github.com/rendley/vegshare/backend/internal/operations/repository"
	"github.com/rendley/vegshare/backend/pkg/config"
	"github.com/rendley/vegshare/backend/pkg/rabbitmq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockOperationsRepository struct {
	mock.Mock
}

func (m *MockOperationsRepository) CreateOperationLog(ctx context.Context, log *operationsModels.OperationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockOperationsRepository) GetOperationLogsForUnit(ctx context.Context, unitID uuid.UUID) ([]operationsModels.OperationLog, error) {
	args := m.Called(ctx, unitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]operationsModels.OperationLog), args.Error(1)
}

func (m *MockOperationsRepository) DeleteOperationLog(ctx context.Context, logID uuid.UUID) error {
	args := m.Called(ctx, logID)
	return args.Error(0)
}

func (m *MockOperationsRepository) UpdateOperationLogStatus(ctx context.Context, logID uuid.UUID, status string) error {
	args := m.Called(ctx, logID, status)
	return args.Error(0)
}

var _ operationsRepository.Repository = &MockOperationsRepository{}

type MockLeasingRepository struct {
	mock.Mock
}

func (m *MockLeasingRepository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leasingModels.Lease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]leasingModels.Lease), args.Error(1)
}

func (m *MockLeasingRepository) GetEnrichedLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]leasingModels.EnrichedLease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]leasingModels.EnrichedLease), args.Error(1)
}

func (m *MockLeasingRepository) CreateLease(ctx context.Context, lease *leasingModels.Lease) error {
	return m.Called(ctx, lease).Error(0)
}

var _ leasingRepository.Repository = &MockLeasingRepository{}

type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) Publish(queueName, body string) error {
	args := m.Called(queueName, body)
	return args.Error(0)
}
func (m *MockRabbitMQClient) Consume(queueName string) (<-chan amqp.Delivery, error) { return nil, nil }
func (m *MockRabbitMQClient) Close()                                                 { m.Called() }

var _ rabbitmq.ClientInterface = &MockRabbitMQClient{}

// --- Tests ---

func TestOperationsService(t *testing.T) {
	ctx := context.Background()
	mockOpsRepo := new(MockOperationsRepository)
	mockLeasingRepo := new(MockLeasingRepository)
	mockRabbitMQ := new(MockRabbitMQClient)
	cfg := &config.Config{
		RabbitMQ: config.RabbitMQConfig{
			Queues: map[string]string{"actions": "actions_queue_test"},
		},
	}

	opsSvc := NewOperationsService(mockOpsRepo, mockLeasingRepo, mockRabbitMQ, cfg)

	t.Run("CreateAction", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			unitID := uuid.New()
			leaseID := uuid.New()
			activeLease := []leasingModels.Lease{{ID: leaseID, UnitID: unitID, UserID: userID, Status: "active", UnitType: "plot"}}
			req := ActionRequest{
				UnitID:     unitID,
				UnitType:   "plot",
				ActionType: "plant",
				Parameters: json.RawMessage(`{"crop_id": "` + uuid.New().String() + `"}`),
			}

			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(activeLease, nil).Once()
			mockOpsRepo.On("CreateOperationLog", ctx, mock.AnythingOfType("*models.OperationLog")).Return(nil).Once()
			mockRabbitMQ.On("Publish", cfg.RabbitMQ.Queues["actions"], mock.AnythingOfType("string")).Return(nil).Once()

			// Act
			logEntry, err := opsSvc.CreateAction(ctx, userID, req)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, logEntry)
			assert.Equal(t, "pending", logEntry.Status)
			mockLeasingRepo.AssertExpectations(t)
			mockOpsRepo.AssertExpectations(t)
			mockRabbitMQ.AssertExpectations(t)
		})

		t.Run("No active lease", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			unitID := uuid.New()
			activeLease := []leasingModels.Lease{}
			req := ActionRequest{
				UnitID:     unitID,
				UnitType:   "plot",
				ActionType: "plant",
			}

			mockLeasingRepo.On("GetLeasesByUserID", ctx, userID).Return(activeLease, nil).Once()

			// Act
			logEntry, err := opsSvc.CreateAction(ctx, userID, req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, logEntry)
			mockLeasingRepo.AssertExpectations(t)
		})
	})
}