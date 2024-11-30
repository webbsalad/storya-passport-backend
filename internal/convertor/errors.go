package convertor

import (
	"errors"

	"github.com/webbsalad/storya-passport-backend/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertError(err error) error {
	switch {
	case errors.Is(err, model.ErrAlreadyExist):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, model.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Errorf(codes.Internal, "server error: %v", err)
}
