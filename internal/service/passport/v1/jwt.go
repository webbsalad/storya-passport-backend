package v1

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func extractTokenClaims(token, secret, expectedType string) (model.Session, error) {
	if secret == "" {
		return model.Session{}, fmt.Errorf("missing jwt secret")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return model.Session{}, model.ErrExpiredToken
		}
		return model.Session{}, fmt.Errorf("invalid token: %w", err)
	}

	if !parsedToken.Valid {
		return model.Session{}, fmt.Errorf("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return model.Session{}, fmt.Errorf("extract token claims")
	}

	if expectedType != "" {
		tokenType, ok := claims["type"].(string)
		if !ok {
			return model.Session{}, fmt.Errorf("type not found: %w", err)
		}

		if tokenType != expectedType {
			return model.Session{}, fmt.Errorf("invalid token type '%s'", expectedType)
		}
	}

	jwtUserID, ok := claims["user_id"].(string)
	if !ok {
		return model.Session{}, fmt.Errorf("user id not found")
	}

	userID, err := model.UserIDFromString(jwtUserID)
	if err != nil {
		return model.Session{}, fmt.Errorf("convert str to user id: %w", err)
	}

	jwtDeviceID, ok := claims["device_id"].(string)
	if !ok {
		return model.Session{}, fmt.Errorf("device id not found")
	}

	deviceID, err := model.DeviceIDFromString(jwtDeviceID)
	if err != nil {
		return model.Session{}, fmt.Errorf("convert str to device id: %w", err)
	}

	jwtVersion, ok := claims["version"].(float64)
	if !ok {
		return model.Session{}, fmt.Errorf("version not found")
	}

	return model.Session{
		UserID:   userID,
		DeviceID: deviceID,
		Version:  int(jwtVersion),
	}, nil

}

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
