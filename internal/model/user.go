package model

import "time"

type User struct {
	Name    string
	EmailId string

	CreatedAt time.Time
	UpdatedAt time.Time
}
