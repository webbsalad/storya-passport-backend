package passport

import (
	"context"

	"github.com/webbsalad/storya-passport-backend/internal/convertor"
	"github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"github.com/webbsalad/storya-passport-backend/internal/utils/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *passport.DeleteRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	userID, err := metadata.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	if err = i.PassportService.Delete(ctx, userID, req.GetEmail()); err != nil {
		return nil, convertor.ConvertError(err)
	}

	return &emptypb.Empty{}, nil

}
