package v1

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func generateTokens(session model.Session, secret string) (model.AuthTokens, error) {
	accessClaims := jwt.MapClaims{
		"user_id":   session.UserID.String(),
		"device_id": session.DeviceID.String(),
		"version":   session.Version,
		"exp":       time.Now().Add(1 * 24 * time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   session.UserID.String(),
		"device_id": session.DeviceID.String(),
		"version":   session.Version,
		"type":      "refresh",
		"exp":       time.Now().Add(365 * 24 * time.Hour).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return model.AuthTokens{}, fmt.Errorf("sign refresh token: %w", err)
	}

	return model.AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
