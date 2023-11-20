package workspace

import (
	channel "github.com/chukwuka-emi/easysync/Channel"
	"gorm.io/gorm"
)

// Workspace ...
type Workspace struct {
	gorm.Model
	Name     string            `json:"name"`
	LogoURL  string            `json:"logoUrl"`
	Channels []channel.Channel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"channels"`
}
