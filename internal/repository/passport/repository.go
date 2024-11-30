package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/model"
)

type Repository interface {
	Register(ctx context.Context, name, passwordHash string) (model.UserID, error)
}
