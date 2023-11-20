package channel

import "gorm.io/gorm"

// Channel Model
type Channel struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        Type   `json:"type"`
	OwnerEmail  string `json:"ownerEmail"`
	WorkspaceID uint   `json:"workspaceId"`
}

// Type ...
type Type string

const (
	//Public ...
	Public Type = "public"
	// Private ...
	Private Type = "private"
)
