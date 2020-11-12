package data

import (
	"os"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewConnection creates the connection to the database
func NewConnection() (*sqlx.DB, error) {
	
	host := os.Getenv("AUTH_DB_HOST")
	port := os.Getenv("AUTH_DB_PORT")
	user := os.Getenv("AUTH_DB_USER")
	DBName := os.Getenv("AUTH_DB_NAME")
	password := os.Getenv("AUTH_DB_PASSWORD")
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, DBName, password)
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}