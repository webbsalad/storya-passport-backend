package metadata

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/model"
	"google.golang.org/grpc/metadata"
)

func GetUserID(ctx context.Context) (model.UserID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return model.UserID{}, fmt.Errorf("missing metadata")
	}

	strUserID := md.Get("user_id")
	if len(strUserID) == 0 {
		return model.UserID{}, fmt.Errorf("missing user id")
	}

	userID, err := model.UserIDFromString(strUserID[0])
	if err != nil {
		return model.UserID{}, fmt.Errorf("convert str to user id: %w", err)
	}

	return userID, nil
}
