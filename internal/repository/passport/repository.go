package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

type Repository interface {
	Register(ctx context.Context, name, passwordHash string) (model.Session, error)
	GetUser(ctx context.Context, userID model.UserID) (model.User, error)
	UpdateUser(ctx context.Context, userID model.UserID, name, passwordHash string) (model.User, error)
	GetPasswordHash(ctx context.Context, name string) (string, error)
	GetSessionInfo(ctx context.Context, name string) (model.Session, error)
	GetTokenVersion(ctx context.Context, userID model.UserID, deviceID model.DeviceID) (int, error)
	UpdateTokenVersion(ctx context.Context, session model.Session) error
	LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error
}
