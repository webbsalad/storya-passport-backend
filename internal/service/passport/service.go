package service

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

type Service interface {
	Register(ctx context.Context, emailID model.EmailID, name, password string) (model.AuthTokens, error)
	GetUser(ctx context.Context, userID model.UserID) (model.User, error)
	UpdateUser(ctx context.Context, userID model.UserID, name, password string) (model.User, error)
	LogIn(ctx context.Context, name, password string) (model.AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.AuthTokens, error)
	LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error
	CheckToken(ctx context.Context, token string) (model.UserID, model.DeviceID, string, error)
}
