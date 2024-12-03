package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) UpdateUser(ctx context.Context, userID model.UserID, name, password string) (model.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.passportRepository.UpdateUser(ctx, userID, name, hashedPassword)
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}
