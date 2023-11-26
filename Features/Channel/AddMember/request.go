package addmember

import "github.com/google/uuid"

type request struct {
	UserID    uuid.UUID `json:"userId" binding:"required"`
	ChannelID uuid.UUID `json:"channelId" binding:"required"`
}
