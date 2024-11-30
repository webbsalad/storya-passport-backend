package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) Register(ctx context.Context, name, passord string) (model.AuthTokens, error) {
	// TODO хеширование пароля, создание токенов
	userID, err := s.passportRepository.Register(ctx, name, passord)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("register user: %w", err)
	}

	return model.NewAuthTokensFromStrings(userID.String(), userID.String()), err
}
