package passport

import (
	"context"
	"fmt"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/model"
	"github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *Implementation) GetUser(ctx context.Context, req *passport.GetUserRequest) (*passport.Passport, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}

	userID, err := model.UserIDFromString(req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("convert str to user id: %w", err)
	}

	user, err := i.PassportService.GetUser(ctx, userID)
	if err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &passport.Passport{
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
