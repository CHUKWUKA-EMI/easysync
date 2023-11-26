package createconversation

import "github.com/google/uuid"

type request struct {
	TargetUserID uuid.UUID `json:"targetUserId" binding:"required"`
}
