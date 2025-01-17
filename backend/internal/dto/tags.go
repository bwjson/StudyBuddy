package dto

type Tag struct {
	ID          uint   `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Title       string `gorm:"size:255;unique;not null" json:"title" binding:"required"`
	Description string `gorm:"size:255;not null" json:"description" binding:"required"`
}

type UserTag struct {
	ID     uint `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	UserID uint `gorm:"not null" json:"user_id"`
	TagID  uint `gorm:"not null" json:"tag_id"`
}
