package invitecollaborator

type request struct {
	CollaboratorEmail string `json:"collaboratorEmail" binding:"required"`
	WorkspaceID       uint   `json:"workspaceID" binding:"required"`
}
