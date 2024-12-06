package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) LogIn(ctx context.Context, name, password string) (model.AuthTokens, error) {
	storedHash, err := s.passportRepository.GetPasswordHash(ctx, name)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get password hash: %w", err)
	}

	if err = comparePassword(storedHash, password); err != nil {
		return model.AuthTokens{}, model.ErrWrongPassword
	}

	session, err := s.passportRepository.GetSessionInfo(ctx, name)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get session info: %w", err)
	}

	tokens, err := generateTokens(session, s.config.JWTSecret)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get tokens: %w", err)
	}

	return tokens, nil

}
