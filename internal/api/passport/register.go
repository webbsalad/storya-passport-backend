package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/model"
	"github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Register(ctx context.Context, req *passport.RegisterRequest) (*passport.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	emailID, err := model.EmailIDFromString(req.GetEmailId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email id: %v", err)
	}

	authTokens, err := i.PassportService.Register(ctx, emailID, req.GetName(), req.GetPassword())
	if err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &passport.RegisterResponse{
		AccessToken:  authTokens.AccessToken,
		RefreshToken: authTokens.RefreshToken,
	}, nil
}
