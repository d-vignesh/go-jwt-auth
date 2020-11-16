package data

import (
	// "os"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/hashicorp/go-hclog"
)

// NewConnection creates the connection to the database
func NewConnection(config *utils.Configurations, logger hclog.Logger) (*sqlx.DB, error) {
	
	// host := os.Getenv("AUTH_DB_HOST")
	// port := os.Getenv("AUTH_DB_PORT")
	// user := os.Getenv("AUTH_DB_USER")
	// DBName := os.Getenv("AUTH_DB_NAME")
	// password := os.Getenv("AUTH_DB_PASSWORD")
	host := config.DBHost 
	port := config.DBPort 
	user := config.DBUser
	dbName := config.DBName
	password := config.DBPass 

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbName, password)
	if config.DATABASE_URL != "" {
		conn = pq.ParseURL(config.DATABASE_URL)
		conn += " sslmode=disable"
	}
	logger.Debug("connection string", conn)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}