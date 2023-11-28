package deletemessage

import "github.com/gocql/gocql"

// MessageDeleteRequest is the payload for chat deletion
type MessageDeleteRequest struct {
	ID             uint64     `json:"id" binding:"required"`
	ConversationID gocql.UUID `json:"conversationId" binding:"required"`
}
