package dto

type JsonInput struct {
	Message string `json:"message"`
}

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"size:255;not null" json:"name"`
	Username     string `gorm:"size:255;unique;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"passwordhash"`
}
