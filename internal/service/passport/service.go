package service

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

type Service interface {
	Register(ctx context.Context, name, passord string) (model.AuthTokens, error)
	Login(ctx context.Context, name, password string) (model.AuthTokens, error)
	LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error
}
