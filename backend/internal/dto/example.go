package dto

type JsonInput struct {
	Message string `json:"message"`
}

type EmailInput struct {
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Name         string `gorm:"size:255;not null" json:"name" binding:"required"`
	Username     string `gorm:"size:255;unique;not null" json:"username" binding:"required"`
	PasswordHash string `gorm:"size:255;not null" json:"password_hash" binding:"required"`
}
