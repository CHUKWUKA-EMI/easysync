package createuser

type createUserRequest struct {
	Email                 string `json:"email" binding:"required"`
	EmailConfirmationCode int    `json:"emailConfirmationCode" binding:"required"`
}
