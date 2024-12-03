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

	ErrWrongPassword     = fmt.Errorf("wrong password: %w", ErrPermissionDenied)
	ErrExpiredToken      = fmt.Errorf("expired token: %w", ErrPermissionDenied)
	ErrExpiredSession    = fmt.Errorf("wrong token version: %w", ErrPermissionDenied)
	ErrWrongTokenType    = fmt.Errorf("wrong token type: %w", ErrPermissionDenied)
	ErrEmailNotConfirmed = fmt.Errorf("email not confirmed: %w", ErrPermissionDenied)

	ErrUserNotFound    = fmt.Errorf("user not found: %w", ErrNotFound)
	ErrDeviceNotFound  = fmt.Errorf("device not found: %w", ErrNotFound)
	ErrMissingMetaData = fmt.Errorf("missing context metadata: %w", ErrNotFound)
)
