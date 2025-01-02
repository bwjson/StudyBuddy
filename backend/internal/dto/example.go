package dto

type JsonInput struct {
	Message string `json:"message"`
}

type EmailInput struct {
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}
