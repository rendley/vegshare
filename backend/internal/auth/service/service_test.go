package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/auth/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) UserExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Check(hashedPassword, inputPassword string) bool {
	args := m.Called(hashedPassword, inputPassword)
	return args.Bool(0)
}

type MockGenerator struct {
	mock.Mock
}

func (m *MockGenerator) GenerateAccessToken(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockGenerator) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// --- Tests ---

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	mockHasher := new(MockHasher)
	mockGenerator := new(MockGenerator)

	authSvc := NewAuthService(mockRepo, mockHasher, mockGenerator)

	t.Run("should register user successfully", func(t *testing.T) {
		// Arrange
		name := "Test User"
		email := "test@example.com"
		password := "password123"
		hashedPassword := "hashed-password"
		refreshToken := "refresh-token"

		mockRepo.On("UserExists", ctx, email).Return(false, nil).Once()
		mockHasher.On("Hash", password).Return(hashedPassword, nil).Once()
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockGenerator.On("GenerateAccessToken", mock.AnythingOfType("uuid.UUID")).Return("access-token", nil).Once()
		mockGenerator.On("GenerateRefreshToken").Return(refreshToken, nil).Once()
		mockRepo.On("SaveRefreshToken", ctx, mock.AnythingOfType("uuid.UUID"), refreshToken).Return(nil).Once()

		// Act
		user, tokens, err := authSvc.Register(ctx, name, email, password)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, tokens)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, hashedPassword, user.PasswordHash)
		assert.Equal(t, "access-token", tokens.AccessToken)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockGenerator.AssertExpectations(t)
	})

	t.Run("should return error if user already exists", func(t *testing.T) {
		// Arrange
		name := "Test User"
		email := "test@example.com"
		password := "password123"

		mockRepo.On("UserExists", ctx, email).Return(true, nil).Once()

		// Act
		user, tokens, err := authSvc.Register(ctx, name, email, password)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, models.ErrUserExists))
		assert.Nil(t, user)
		assert.Nil(t, tokens)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	mockHasher := new(MockHasher)
	mockGenerator := new(MockGenerator)

	authSvc := NewAuthService(mockRepo, mockHasher, mockGenerator)

	t.Run("should login user successfully", func(t *testing.T) {
		// Arrange
		email := "test@example.com"
		password := "password123"
		hashedPassword := "hashed-password"
		user := &models.User{ID: uuid.New(), Email: email, PasswordHash: hashedPassword}
		refreshToken := "refresh-token"

		mockRepo.On("GetUserByEmail", ctx, email).Return(user, nil).Once()
		mockHasher.On("Check", hashedPassword, password).Return(true).Once()
		mockGenerator.On("GenerateAccessToken", user.ID).Return("access-token", nil).Once()
		mockGenerator.On("GenerateRefreshToken").Return(refreshToken, nil).Once()
		mockRepo.On("SaveRefreshToken", ctx, user.ID, refreshToken).Return(nil).Once()

		// Act
		loggedInUser, tokens, err := authSvc.Login(ctx, email, password)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, loggedInUser)
		assert.NotNil(t, tokens)
		assert.Equal(t, user.ID, loggedInUser.ID)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockGenerator.AssertExpectations(t)
	})

	t.Run("should return error for invalid credentials on wrong email", func(t *testing.T) {
		// Arrange
		email := "wrong@example.com"
		password := "password123"

		mockRepo.On("GetUserByEmail", ctx, email).Return(nil, models.ErrInvalidCredentials).Once()

		// Act
		loggedInUser, tokens, err := authSvc.Login(ctx, email, password)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, models.ErrInvalidCredentials))
		assert.Nil(t, loggedInUser)
		assert.Nil(t, tokens)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid credentials on wrong password", func(t *testing.T) {
		// Arrange
		email := "test@example.com"
		password := "wrong-password"
		hashedPassword := "hashed-password"
		user := &models.User{ID: uuid.New(), Email: email, PasswordHash: hashedPassword}

		mockRepo.On("GetUserByEmail", ctx, email).Return(user, nil).Once()
		mockHasher.On("Check", hashedPassword, password).Return(false).Once()

		// Act
		loggedInUser, tokens, err := authSvc.Login(ctx, email, password)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, models.ErrInvalidCredentials))
		assert.Nil(t, loggedInUser)
		assert.Nil(t, tokens)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
