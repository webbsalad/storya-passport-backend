package model

import "time"

type User struct {
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}
