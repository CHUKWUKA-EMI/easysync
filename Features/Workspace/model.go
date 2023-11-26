package workspace

import (
	data "github.com/chukwuka-emi/easysync/Data"
	channel "github.com/chukwuka-emi/easysync/Features/Channel"
	"github.com/google/uuid"
)

// Workspace ...
type Workspace struct {
	data.BaseModel
	Name     string            `json:"name"`
	LogoURL  string            `json:"logoUrl"`
	Channels []channel.Channel `gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"channels"`
}

// Invite ...
type Invite struct {
	data.BaseModel
	Token       string    `json:"token"`
	Status      string    `gorm:"type:enum('invite_pending','invite_accepted');" json:"status"`
	WorkspaceID uuid.UUID `gorm:"type:char(36);index:idx_member,priority:12;" json:"workspaceId"`
	UserID      uuid.UUID `gorm:"type:char(36);index:idx_member;" json:"userId"`
}

const (
	// InvitePending is a constant for pending workspace invites
	InvitePending string = "invite_pending"
	// InviteAccepted is a constant for accepted workspace invites
	InviteAccepted string = "invite_accepted"
)
