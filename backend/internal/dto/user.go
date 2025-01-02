package dto

type UsersWithPagination struct {
	User       []User `json:"users"`
	TotalCount int    `json:"totalCount"`
}

type User struct {
	ID                uint   `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Name              string `gorm:"size:255;not null" json:"name" binding:"required"`
	Email             string `gorm:"size:255;unique;not null" json:"email" binding:"required"`
	Username          string `gorm:"size:255;unique;not null" json:"username" binding:"required"`
	PasswordHash      string `gorm:"size:255;not null" json:"password_hash" binding:"required"`
	IsActive          bool   `gorm:"default:false" json:"is_active" swaggerignore:"true"`
	VerificationToken string `gorm:"size:255;" json:"verification_token" swaggerignore:"true"`
}
