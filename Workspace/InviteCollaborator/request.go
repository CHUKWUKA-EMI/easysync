package invitecollaborator

import uuid "github.com/satori/go.uuid"

type request struct {
	CollaboratorEmail string    `json:"collaboratorEmail" binding:"required"`
	WorkspaceID       uuid.UUID `json:"workspaceID" binding:"required"`
}
