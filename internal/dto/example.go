package dto

type JsonInput struct {
	Message string `json:"message"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:255;not null"`
	Username     string `gorm:"size:255;unique;not null"`
	PasswordHash string `gorm:"size:255;not null"`
}
