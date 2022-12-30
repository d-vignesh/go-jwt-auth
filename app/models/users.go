package models

import "time"

type User struct {
	ID         string    `json:"id" db:"id"`
	Email      string    `json:"email" validate:"required" db:"email"`
	Password   string    `json:"password" validate:"required" db:"password"`
	Username   string    `json:"username" db:"username"`
	Token      string    `json:"token" db:"token"`
	IsVerified bool      `json:"isVerified" db:"is_verified"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}
