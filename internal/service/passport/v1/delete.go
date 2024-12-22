package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) Delete(ctx context.Context, userID model.UserID, emailID model.EmailID) error {
	if err := s.passportRepository.Delete(ctx, userID, emailID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
