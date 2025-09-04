package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rendley/vegshare/backend/internal/leasing/domain"
	"github.com/rendley/vegshare/backend/internal/leasing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockLeasingRepository struct {
	mock.Mock
}

func (m *MockLeasingRepository) CreateLease(ctx context.Context, lease *models.Lease) error {
	args := m.Called(ctx, lease)
	return args.Error(0)
}

func (m *MockLeasingRepository) GetEnrichedLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.EnrichedLease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.EnrichedLease), args.Error(1)
}

func (m *MockLeasingRepository) GetLeasesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Lease, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Lease), args.Error(1)
}

type MockUnitManager struct {
	mock.Mock
}

func (m *MockUnitManager) GetLeasableUnit(ctx context.Context, unitID uuid.UUID) (domain.LeasableUnit, error) {
	args := m.Called(ctx, unitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.LeasableUnit), args.Error(1)
}

func (m *MockUnitManager) UpdateUnitStatus(ctx context.Context, unitID uuid.UUID, status string) error {
	args := m.Called(ctx, unitID, status)
	return args.Error(0)
}

func (m *MockUnitManager) WithTx(tx *sqlx.Tx) domain.UnitManager {
	return m
}

// MockLeasableUnit - это мок для самого юнита, чтобы мы могли контролировать его статус.
type MockLeasableUnit struct {
	mock.Mock
}

func (m *MockLeasableUnit) GetID() uuid.UUID    { return uuid.New() }
func (m *MockLeasableUnit) GetUnitType() string { return string(models.UnitTypePlot) }
func (m *MockLeasableUnit) GetStatus() string {
	args := m.Called()
	return args.String(0)
}

// --- Tests ---

func TestLeasingService(t *testing.T) {
	ctx := context.Background()
	mockLeasingRepo := new(MockLeasingRepository)
	mockUnitMgr := new(MockUnitManager)

	leasingSvc := NewLeasingService(nil, mockLeasingRepo)
	leasingSvc.RegisterUnitManager(models.UnitTypePlot, mockUnitMgr)

	t.Run("CreateLease", func(t *testing.T) {
		t.Run("Unit not found", func(t *testing.T) {
			userID := uuid.New()
			unitID := uuid.New()
			mockUnitMgr.On("GetLeasableUnit", ctx, unitID).Return(nil, errors.New("not found")).Once()

			lease, err := leasingSvc.CreateLease(ctx, userID, unitID, models.UnitTypePlot)

			assert.Error(t, err)
			assert.Nil(t, lease)
			mockUnitMgr.AssertExpectations(t)
		})

		t.Run("Unit not available", func(t *testing.T) {
			userID := uuid.New()
			unitID := uuid.New()

			mockUnit := new(MockLeasableUnit)
			mockUnit.On("GetStatus").Return("rented")

			mockUnitMgr.On("GetLeasableUnit", ctx, unitID).Return(mockUnit, nil).Once()

			lease, err := leasingSvc.CreateLease(ctx, userID, unitID, models.UnitTypePlot)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "недоступен для аренды")
			assert.Nil(t, lease)
			mockUnitMgr.AssertExpectations(t)
		})

		t.Run("Unit Manager not registered", func(t *testing.T) {
			userID := uuid.New()
			unitID := uuid.New()

			// Пытаемся арендовать несуществующий тип юнита
			lease, err := leasingSvc.CreateLease(ctx, userID, unitID, "non_existent_type")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "не зарегистрирован")
			assert.Nil(t, lease)
		})
	})

	t.Run("GetMyEnrichedLeases", func(t *testing.T) {
		userID := uuid.New()
		expectedLeases := []models.EnrichedLease{{Lease: models.Lease{ID: uuid.New(), UserID: userID}}}
		mockLeasingRepo.On("GetEnrichedLeasesByUserID", ctx, userID).Return(expectedLeases, nil).Once()

		leases, err := leasingSvc.GetMyEnrichedLeases(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedLeases, leases)
		mockLeasingRepo.AssertExpectations(t)
	})
}
