package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/user/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *models.UserProfile) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// --- Tests ---

func TestUserService(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	userSvc := NewUserService(mockRepo)

	t.Run("GetUser", func(t *testing.T) {
		t.Run("should get user successfully", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			expectedUser := &models.UserProfile{ID: userID, Name: "Test"}
			mockRepo.On("GetUserByID", ctx, userID).Return(expectedUser, nil).Once()

			// Act
			user, err := userSvc.GetUser(ctx, userID)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expectedUser, user)
			mockRepo.AssertExpectations(t)
		})

		t.Run("should return error when user not found", func(t *testing.T) {
			// Arrange
			userID := uuid.New()
			mockRepo.On("GetUserByID", ctx, userID).Return(nil, errors.New("not found")).Once()

			// Act
			user, err := userSvc.GetUser(ctx, userID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, user)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Arrange
		userID := uuid.New()
		originalUser := &models.UserProfile{ID: userID, Name: "Original Name"}
		newName := "New Name"
		newAvatar := "http://example.com/avatar.png"

		mockRepo.On("GetUserByID", ctx, userID).Return(originalUser, nil).Once()
		mockRepo.On("UpdateUser", ctx, mock.AnythingOfType("*models.UserProfile")).Return(nil).Once()

		// Act
		updatedUser, err := userSvc.UpdateUser(ctx, userID, newName, newAvatar)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, newName, updatedUser.Name)
		assert.Equal(t, &newAvatar, updatedUser.AvatarURL)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// Arrange
		userID := uuid.New()
		mockRepo.On("DeleteUser", ctx, userID).Return(nil).Once()

		// Act
		err := userSvc.DeleteUser(ctx, userID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
