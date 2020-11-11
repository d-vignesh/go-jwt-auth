package handlers

import (
	"fmt"

	"github.com/d-vignesh/go-jwt-auth/data"
	"github.com/d-vignesh/go-jwt-auth/service"
	"github.com/hashicorp/go-hclog"
)

// UserKey is used as a key for storing the User object in context at middleware
type UserKey struct{}

type UserIDKey struct{}

// UserHandler provides means to perform operations on user object
type UserHandler struct {
	logger 		hclog.Logger
	validator	*data.Validation
	repo 		data.Repository
	authService service.Authentication
}

func NewUserHandler(l hclog.Logger, v *data.Validation, r data.Repository, authService service.Authentication) *UserHandler {
	return &UserHandler {l, v, r, authService}
}

// GenericError is a generic error message returned by server
type GenericError struct {
	Message  string  `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages  []string	`json:"messages"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken string `json:"access_token"`
}

var ErrUserAlreadyExists = fmt.Sprintf("user already exists with the given email")
var ErrUserNotFound = fmt.Sprintf("invalid email or password")

var PgDuplicateKeyMsg = "duplicate key value violates unique constraint"
var PgNoRowsMsg = "no rows in result set"