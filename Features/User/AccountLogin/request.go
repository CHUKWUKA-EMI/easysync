package accountlogin

type accountLoginRequest struct {
	Email     string `json:"email" binding:"required"`
	LoginCode int    `json:"loginCode" binding:"required"`
}
