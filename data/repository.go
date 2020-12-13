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
	StoreVerificationData(ctx context.Context, verificationData *VerificationData) error
	GetVerificationData(ctx context.Context, email string, verificationDataType VerificationDataType) (*VerificationData, error)
	UpdateUserVerificationStatus(ctx context.Context, email string, status bool) error
	DeleteVerificationData(ctx context.Context, email string, verificationDataType VerificationDataType) error
	UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error
}
