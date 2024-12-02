package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) CheckToken(ctx context.Context, token string) (model.UserID, model.DeviceID, string, error) {
	session, tokenType, err := extractTokenClaims(token, s.config.JWTSecret)
	if err != nil {
		return model.UserID{}, model.DeviceID{}, "", fmt.Errorf("extract claims: %w", err)
	}

	storedVersion, err := s.passportRepository.GetTokenVersion(ctx, session.UserID, session.DeviceID)
	if err != nil {
		return model.UserID{}, model.DeviceID{}, "", fmt.Errorf("get stored version: %w", err)
	}

	if storedVersion != session.Version {
		return model.UserID{}, model.DeviceID{}, "", model.ErrWrongPassword
	}

	return session.UserID, session.DeviceID, tokenType, nil
}
