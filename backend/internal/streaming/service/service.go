package service

import (
	"context"

	"github.com/google/uuid"
	cameraService "github.com/rendley/vegshare/backend/internal/camera/service"
)

// Service defines the contract for the streaming service.
type Service interface {
	// We can add methods here later, for example, to authorize the stream
	AuthorizeStream(ctx context.Context, cameraID uuid.UUID) (bool, error)
}

// service implements the Service interface.
type service struct {
	cameraSvc cameraService.Service
}

// NewService is a constructor for the streaming service.
func NewService(cameraSvc cameraService.Service) Service {
	return &service{cameraSvc: cameraSvc}
}

func (s *service) AuthorizeStream(ctx context.Context, cameraID uuid.UUID) (bool, error) {
	// For now, we just check if the camera exists.
	// In the future, we could check if the user has rights to view this camera (e.g., if they lease the plot).
	_, err := s.cameraSvc.GetCameraByID(ctx, cameraID)
	if err != nil {
		return false, err
	}
	return true, nil
}
