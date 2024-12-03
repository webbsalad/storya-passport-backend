package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) GetUser(ctx context.Context, userID model.UserID) (model.User, error) {
	user, err := s.passportRepository.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
