package chat

import (
	"time"

	"github.com/gocql/gocql"
)

// Conversation is the model for Easysync conversations table
type Conversation struct {
	ID        gocql.UUID   `json:"id"`
	Members   []gocql.UUID `json:"members"`
	IsOpen    bool         `json:"isOpen"`
	CreatedAt time.Time    `json:"createdAt"`
}

// Chat is the model for Easysync chats table
type Chat struct {
	ID             uint64     `json:"id"`
	ConversationID gocql.UUID `json:"conversationId"`
	SenderID       gocql.UUID `json:"senderId"`
	Content        string     `json:"content"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
