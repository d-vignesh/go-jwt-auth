package data

import (
	"context"
)

// Repository is an interface for the storage implementation of the auth service
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	UpdateUsername(ctx context.Context, user *User) error
}
