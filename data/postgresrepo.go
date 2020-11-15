package data

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/hashicorp/go-hclog"
)

type PostgresRepository struct {
	db *sqlx.DB 
	logger hclog.Logger 
}

func NewPostgresRepository(db *sqlx.DB, logger hclog.Logger) *PostgresRepository {
	return &PostgresRepository{db, logger}
}

func (repo *PostgresRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewV4().String()
	repo.logger.Info("creating user", hclog.Fmt("%#v", user))
	query := "insert into users (id, email, username, password) values ($1, $2, $3, $4)"
	_, err := repo.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	repo.logger.Debug("querying for user with email", email)
	query := "select * from users where email = $1"
	var user User
	if err := repo.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	repo.logger.Debug("querying for user with id", userID)
	query := "select * from users where id = $1"
	var user User
	if err := repo.db.GetContext(ctx, &user, query, userID); err != nil {
		return nil, err
	}
	return &user, nil
}