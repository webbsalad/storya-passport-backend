package v1

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (model.AuthTokens, error) {
	session, tokenType, err := extractTokenClaims(refreshToken, s.config.JWTSecret)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("extract session info: %w", err)
	}

	if tokenType != "refresh" {
		return model.AuthTokens{}, model.ErrWrongTokenType
	}

	storedVersion, err := s.passportRepository.GetTokenVersion(ctx, session.UserID, session.DeviceID)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get stored version: %w", err)
	}

	if storedVersion != session.Version {
		return model.AuthTokens{}, model.ErrExpiredSession
	}

	session.Version += 1
	authTokens, err := generateTokens(session, s.config.JWTSecret)
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("get auth tokens: %w", err)
	}

	if err = s.passportRepository.UpdateTokenVersion(ctx, session); err != nil {
		return model.AuthTokens{}, fmt.Errorf("store token version: %w", err)
	}

	return authTokens, nil

}
