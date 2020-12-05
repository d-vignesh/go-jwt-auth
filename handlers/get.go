package handlers

import (
	"net/http"

	"github.com/d-vignesh/go-jwt-auth/data"
)

// RefreshToken handles refresh token request
func (uh *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(UserKey{}).(data.User)
	accessToken, err := uh.authService.GenerateAccessToken(&user)
	if err != nil {
		uh.logger.Error("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	data.ToJSON(&TokenResponse{AccessToken: accessToken}, w)
}

// Greet request greet request
func (uh *UserHandler) Greet(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(UserIDKey{}).(string)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, " + userID))
}
