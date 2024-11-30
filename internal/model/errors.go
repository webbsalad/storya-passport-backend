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
)
