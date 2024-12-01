package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) RefreshToken(ctx context.Context, req *passport.RefreshTokenRequest) (*passport.RefreshTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	authTokens, err := i.PassportService.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &passport.RefreshTokenResponse{
		AccessToken:  authTokens.AccessToken,
		RefreshToken: authTokens.RefreshToken,
	}, nil
}
