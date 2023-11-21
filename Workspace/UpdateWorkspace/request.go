package updateworkspace

import uuid "github.com/satori/go.uuid"

type request struct {
	ID   uuid.UUID `json:"id" binding:"required"`
	Name string    `json:"name" binding:"required"`
}
