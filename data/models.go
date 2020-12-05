package data

import (
	"time"
)

// User is the data type for user object
type User struct {
	ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" sql:"username"`
	TokenHash string    `json:"tokenhash" sql:"tokenhash"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
}
