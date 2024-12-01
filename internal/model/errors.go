package model

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExist     = errors.New("already exist")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotFound         = errors.New("not found")
)

var (
	ErrUserAlreadyExist = fmt.Errorf("user already exist: %w", ErrAlreadyExist)

	ErrWrongPassword = fmt.Errorf("wrong password: %w", ErrPermissionDenied)

	ErrDeviceNotFound  = fmt.Errorf("device not found: %w", ErrNotFound)
	ErrMissingMetaData = fmt.Errorf("missing context metadata: %w", ErrNotFound)
)
