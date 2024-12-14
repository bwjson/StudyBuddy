package dto

type JsonInput struct {
	Message string `json:"message"`
}

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:255;not null" json:"name" binding:"required"`
	Username     string `gorm:"size:255;unique;not null" json:"username" binding:"required"`
	PasswordHash string `gorm:"size:255;not null" json:"password_hash" binding:"required"`
}
