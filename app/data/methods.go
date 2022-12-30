package data

import (
	"context"
	"fmt"
	"time"

	"github.com/caseyrwebb/go-jwt-auth/app/models"
	uuid "github.com/google/uuid"
)

// Create inserts the given user into the database
func (d *DB) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	d.logger.Info(fmt.Sprintf("%s %v", "creating user", user))
	query := "insert into users (id, email, username, password, token, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err := d.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Password, user.Token, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves the user object having the given email, else returns error
func (d *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	d.logger.Debug(fmt.Sprintf("%s %v", "querying for user with email", email))
	query := "select * from users where email = $1"
	var user models.User
	if err := d.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	d.logger.Debug(fmt.Sprintf("%s %v", "read users", user))
	return &user, nil
}

// GetUserByID retrieves the user object having the given ID, else returns error
func (d *DB) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	d.logger.Debug(fmt.Sprintf("%s %v", "querying for user with id", userID))
	query := "select * from users where id = $1"
	var user models.User
	if err := d.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUsername updates the username of the given user
func (d *DB) UpdateUsername(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := "update users set username = $1, updatedat = $2 where id = $3"
	if _, err := d.db.ExecContext(ctx, query, user.Username, user.UpdatedAt, user.ID); err != nil {
		return err
	}
	return nil
}

// UpdateUserVerificationStatus updates user verification status to true
func (d *DB) UpdateUserVerificationStatus(ctx context.Context, email string, status bool) error {

	query := "update users set isverified = $1 where email = $2"
	if _, err := d.db.ExecContext(ctx, query, status, email); err != nil {
		return err
	}
	return nil
}

// StoreMailVerificationData adds a mail verification data to db
func (d *DB) StoreVerificationData(ctx context.Context, verificationData *models.VerificationData) error {

	query := "insert into verifications(email, code, expiresat, type) values($1, $2, $3, $4)"
	_, err := d.db.ExecContext(ctx, query, verificationData.Email, verificationData.Code, verificationData.ExpiresAt, verificationData.Type)
	return err
}

// GetMailVerificationCode retrieves the stored verification code.
func (d *DB) GetVerificationData(ctx context.Context, email string, verificationDataType models.VerificationDataType) (*models.VerificationData, error) {

	query := "select * from verifications where email = $1 and type = $2"

	var verificationData models.VerificationData
	if err := d.db.GetContext(ctx, &verificationData, query, email, verificationDataType); err != nil {
		return nil, err
	}
	return &verificationData, nil
}

// DeleteMailVerificationData deletes a used verification data
func (d *DB) DeleteVerificationData(ctx context.Context, email string, verificationDataType models.VerificationDataType) error {

	query := "delete from verifications where email = $1 and type = $2"
	_, err := d.db.ExecContext(ctx, query, email, verificationDataType)
	return err
}

// UpdatePassword updates the user password
func (d *DB) UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error {

	query := "update users set password = $1, token = $2 where id = $3"
	_, err := d.db.ExecContext(ctx, query, password, tokenHash, userID)
	return err
}
