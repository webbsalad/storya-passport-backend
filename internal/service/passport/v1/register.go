package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) Register(ctx context.Context, emailID model.EmailID, name, password string) (model.AuthTokens, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("hashing: %w", err)
	}

	session, err := s.passportRepository.Register(ctx, emailID, name, hashedPassword)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("register user: %w", err)
	}

	tokens, err := generateTokens(session, s.config.JWTSecret)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get tokens: %w", err)
	}

	return tokens, nil

}
