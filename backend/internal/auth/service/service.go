
package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/auth/models"
	"github.com/rendley/vegshare/backend/internal/auth/repository"
	"github.com/rendley/vegshare/backend/pkg/jwt"
	"github.com/rendley/vegshare/backend/pkg/security"
)

// AuthService defines the interface for authentication service.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, *models.TokenPair, error)
	Login(ctx context.Context, email, password string) (*models.User, *models.TokenPair, error)
}

// authService is the implementation of AuthService.
type authService struct {
	repo           repository.AuthRepository
	passwordHasher security.PasswordHasher
	jwtGenerator   jwt.Generator
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo repository.AuthRepository, hasher security.PasswordHasher, jwtGen jwt.Generator) AuthService {
	return &authService{
		repo:           repo,
		passwordHasher: hasher,
		jwtGenerator:   jwtGen,
	}
}

// Register handles user registration.
func (s *authService) Register(ctx context.Context, email, password string) (*models.User, *models.TokenPair, error) {
	exists, err := s.repo.UserExists(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, models.ErrUserExists
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// Login handles user login.
func (s *authService) Login(ctx context.Context, email, password string) (*models.User, *models.TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, models.ErrInvalidCredentials
	}

	if !s.passwordHasher.Check(user.PasswordHash, password) {
		return nil, nil, models.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken); err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// generateTokens generates new access and refresh tokens.
func (s *authService) generateTokens(userID uuid.UUID) (*models.TokenPair, error) {
	accessToken, err := s.jwtGenerator.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtGenerator.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
