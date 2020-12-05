package data

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

// PostgresRepository has the implementation of the db methods.
type PostgresRepository struct {
	db     *sqlx.DB
	logger hclog.Logger
}

// NewPostgresRepository returns a new PostgresRepository instance
func NewPostgresRepository(db *sqlx.DB, logger hclog.Logger) *PostgresRepository {
	return &PostgresRepository{db, logger}
}

// Create inserts the given user into the database
func (repo *PostgresRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.logger.Info("creating user", hclog.Fmt("%#v", user))
	query := "insert into users (id, email, username, password, tokenhash, createdat, updatedat) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err := repo.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Password, user.TokenHash, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves the user object having the given email, else returns error
func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	repo.logger.Debug("querying for user with email", email)
	query := "select * from users where email = $1"
	var user User
	if err := repo.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	repo.logger.Debug("read users", hclog.Fmt("%#v", user))
	return &user, nil
}

// GetUserByID retrieves the user object having the given ID, else returns error
func (repo *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	repo.logger.Debug("querying for user with id", userID)
	query := "select * from users where id = $1"
	var user User
	if err := repo.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUsername updates the username of the given user
func (repo *PostgresRepository) UpdateUsername(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	query := "update users set username = $1, updated_at = $2 where id = $3"
	if _, err := repo.db.ExecContext(ctx, query, user.Username, user.UpdatedAt, user.ID); err != nil {
		return err
	}
	return nil
}
