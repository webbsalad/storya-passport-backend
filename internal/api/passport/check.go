package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CheckToken(ctx context.Context, req *passport.CheckTokenRequest) (*passport.CheckTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, deviceID, tokenType, err := i.PassportService.CheckToken(ctx, req.GetToken())
	if err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &passport.CheckTokenResponse{
		UserId:    userID.String(),
		DeviceId:  deviceID.String(),
		TokenType: tokenType,
	}, nil
}
