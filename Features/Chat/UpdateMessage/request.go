package updatemessage

import (
	"time"

	"github.com/gocql/gocql"
)

// MessageUpdateRequest is the payload for chat update
type MessageUpdateRequest struct {
	ID             uint64     `json:"id" binding:"required"`
	ConversationID gocql.UUID `json:"conversationId" binding:"required"`
	Content        string     `json:"content" binding:"required"`
	UpdatedAt      time.Time  `json:"updatedAt" binding:"required"`
}
