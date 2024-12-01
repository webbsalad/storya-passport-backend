package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error {
	if err := s.passportRepository.LogOut(ctx, userID, deviceID); err != nil {
		return fmt.Errorf("log out: %w", err)
	}

	return nil
}
