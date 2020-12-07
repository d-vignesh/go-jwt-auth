package handlers

import (
	"fmt"

	"github.com/d-vignesh/go-jwt-auth/data"
	"github.com/d-vignesh/go-jwt-auth/service"
	"github.com/hashicorp/go-hclog"
)

// UserKey is used as a key for storing the User object in context at middleware
type UserKey struct{}

// UserIDKey is used as a key for storing the UserID in context at middleware
type UserIDKey struct{}

// UserHandler wraps instances needed to perform operations on user object
type UserHandler struct {
	logger      hclog.Logger
	validator   *data.Validation
	repo        data.Repository
	authService service.Authentication
}

// NewUserHandler returns a new UserHandler instance
func NewUserHandler(l hclog.Logger, v *data.Validation, r data.Repository, authService service.Authentication) *UserHandler {
	return &UserHandler{l, v, r, authService}
}

// // GenericError wraps any generic error returned by server
// type GenericError struct {
// 	Error string `json:"error"`
// }

// // GenericMessage wraps any generic message returned by server
// type GenericMessage struct {
// 	Message string `json:"message"`
// }

// GenericResponse is the format of our response
type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Errors []string `json:"errors"`
}

// Below data types are used for encoding and decoding b/t go types and json
type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Username     string `json:"username"`
}

type UsernameUpdate struct {
	Username string `json:"username"`
}

var ErrUserAlreadyExists = fmt.Sprintf("user already exists with the given email")
var ErrUserNotFound = fmt.Sprintf("no user account exists with given email")

var PgDuplicateKeyMsg = "duplicate key value violates unique constraint"
var PgNoRowsMsg = "no rows in result set"
