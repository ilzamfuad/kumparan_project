package model

import "github.com/google/uuid"

type Author struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
