package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/d-vignesh/go-jwt-auth/data"
	"github.com/d-vignesh/go-jwt-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

// Signup handles signup request
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey{}).(data.User)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		uh.logger.Error("unable to hash password", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}

	user.Password = string(hashedPass)
	user.TokenHash = utils.GenerateRandomString(15)

	err = uh.repo.Create(context.Background(), &user)
	if err != nil {
		uh.logger.Error("unable to insert user to database", "error", err)
		errMsg := err.Error()
		if strings.Contains(errMsg, PgDuplicateKeyMsg) {
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Error: ErrUserAlreadyExists}, w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&GenericError{Error: errMsg}, w)
		}
		return
	}

	uh.logger.Debug("User created successfully")
	w.WriteHeader(http.StatusCreated)
	data.ToJSON(&GenericMessage{Message: "user created successfully"}, w)
}

// Login handles login request
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	reqUser := r.Context().Value(UserKey{}).(data.User)

	user, err := uh.repo.GetUserByEmail(context.Background(), reqUser.Email)
	if err != nil {
		uh.logger.Error("error fetching the user", "error", err)
		errMsg := err.Error()
		if strings.Contains(errMsg, PgNoRowsMsg) {
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Error: ErrUserNotFound}, w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&GenericError{Error: err.Error()}, w)
		}
		return
	}

	if valid := uh.authService.Authenticate(&reqUser, user); !valid {
		uh.logger.Debug("Authetication of user failed")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		data.ToJSON(&GenericError{Error: "incorrect password"}, w)
		return
	}

	accessToken, err := uh.authService.GenerateAccessToken(user)
	if err != nil {
		uh.logger.Error("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}
	refreshToken, err := uh.authService.GenerateRefreshToken(user)
	if err != nil {
		uh.logger.Error("unable to generate refresh token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}

	uh.logger.Debug("successfully generated token", "accesstoken", accessToken, "refreshtoken", refreshToken)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	data.ToJSON(&AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken, Username: user.Username}, w)
}
