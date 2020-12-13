package handlers

import (
	"context"
	"net/http"

	"github.com/d-vignesh/go-jwt-auth/data"
)

// UpdateUsername handles username update request
func (ah *AuthHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	user := &data.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
		ah.logger.Error("unable to decode user json", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		// data.ToJSON(&GenericError{Error: err.Error()}, w)
		data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
		return
	}

	user.ID = r.Context().Value(UserIDKey{}).(string)
	ah.logger.Debug("udpating username for user : ", user)

	err = ah.repo.UpdateUsername(context.Background(), user)
	if err != nil {
		ah.logger.Error("unable to update username", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		// data.ToJSON(&GenericError{Error: err.Error()}, w)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to update username. Please try again later"}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	// data.ToJSON(&UsernameUpdate{Username: user.Username}, w)
	data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "Successfully updated username",
		Data:    &UsernameUpdate{Username: user.Username},
	}, w)
}
