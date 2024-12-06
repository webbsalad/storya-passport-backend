package model

import "time"

type User struct {
	Name    string
	EmailId EmailID

	CreatedAt time.Time
	UpdatedAt time.Time
}
