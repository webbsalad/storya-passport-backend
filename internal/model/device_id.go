package model

import (
	"fmt"

	"github.com/google/uuid"
)

type DeviceID uuid.UUID

func (id DeviceID) String() string {
	return uuid.UUID(id).String()
}

func DeviceIDFromString(s string) (DeviceID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return DeviceID{}, fmt.Errorf("parse: %w", err)
	}
	return DeviceID(id), nil
}
