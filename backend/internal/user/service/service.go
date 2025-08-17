package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rendley/vegshare/backend/internal/user/models"
	"github.com/rendley/vegshare/backend/internal/user/repository"
)

// UserService defines the interface for user service.
type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
	UpdateUser(ctx context.Context, id uuid.UUID, fullName, avatarURL string) (*models.UserProfile, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// userService is the implementation of UserService.
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetUser retrieves a user by their ID.
func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	return s.repo.GetUserByID(ctx, id)
}

// UpdateUser updates a user's profile.
func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, fullName, avatarURL string) (*models.UserProfile, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = fullName
	user.AvatarURL = avatarURL

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by their ID.
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}