package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

type Repository interface {
	Register(ctx context.Context, name, passwordHash string) (model.Session, error)
	GetPasswordHash(ctx context.Context, name string) (string, error)
	GetSessionInfo(ctx context.Context, name string) (model.Session, error)
	LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error
}
