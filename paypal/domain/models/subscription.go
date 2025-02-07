package models

type Subscription struct {
	Email      string `json:"email"`
	CardNumber string `json:"card_number" db:"card_number"`
}
