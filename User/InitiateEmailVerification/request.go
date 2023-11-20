package initiateemailverification

type request struct{
	Email string `json:"email" binding:"required"`
}