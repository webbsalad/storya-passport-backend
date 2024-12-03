package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/utils/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) LogOut(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := metadata.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	deviceID, err := metadata.GetDeviceID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	if err = i.PassportService.LogOut(ctx, userID, deviceID); err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &emptypb.Empty{}, nil

}
