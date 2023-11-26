package channel

import (
	data "github.com/chukwuka-emi/easysync/Data"
	"github.com/google/uuid"
)

// Channel Model
type Channel struct {
	// gorm.Model
	data.BaseModel
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        Type      `json:"type"`
	OwnerEmail  string    `json:"ownerEmail"`
	WorkspaceID uuid.UUID `gorm:"type:char(36);" json:"workspaceId"`
}

// Type ...
type Type string

const (
	//Public ...
	Public Type = "public"
	// Private ...
	Private Type = "private"
)
