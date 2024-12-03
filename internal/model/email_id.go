package model

import (
	"fmt"

	"github.com/google/uuid"
)

type EmailID uuid.UUID

func (id EmailID) String() string {
	return uuid.UUID(id).String()
}

func EmailIDFromString(s string) (EmailID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EmailID{}, fmt.Errorf("parse: %w", err)
	}
	return EmailID(id), nil
}
