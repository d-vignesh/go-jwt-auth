package handlers

import (
	"context"
	"net/http"

	"github.com/d-vignesh/go-jwt-auth/data"
)

// UpdateUsername handles username update request
func (uh *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	user := &data.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
		uh.logger.Error("unable to decode user json", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}

	user.ID = r.Context().Value(UserIDKey{}).(string)
	uh.logger.Debug("udpating username for user : ", user)

	err = uh.repo.UpdateUsername(context.Background(), user)
	if err != nil {
		uh.logger.Error("unable to update username", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Error: err.Error()}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	data.ToJSON(&UsernameUpdate{Username: user.Username}, w)
}
