package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caseyrwebb/go-jwt-auth/app/data"
	"github.com/caseyrwebb/go-jwt-auth/app/models"
	"github.com/caseyrwebb/go-jwt-auth/app/utils"
)

// UpdateUsername handles username update request
func (ah *AuthHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	user := &models.User{}
	err := data.FromJSON(user, r.Body)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %v", "unable to decode user json", "error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		// data.ToJSON(&GenericError{Error: err.Error()}, w)d
		data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
		return
	}

	user.ID = r.Context().Value(UserIDKey{}).(string)
	ah.logger.Debug(fmt.Sprintf("%s %v", "udpating username for user : ", user))

	err = ah.db.UpdateUsername(context.Background(), user)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to update username", "error", err))
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

// PasswordReset handles the password reset request
func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	passResetReq := &PasswordResetReq{}
	err := data.FromJSON(passResetReq, r.Body)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %v", "unable to decode password reset request json", "error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
		return
	}

	userID := r.Context().Value(UserIDKey{}).(string)
	user, err := ah.db.GetUserByID(context.Background(), userID)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to retrieve the user from db", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to reset password. Please try again later"}, w)
		return
	}

	verificationData, err := ah.db.GetVerificationData(context.Background(), user.Email, models.PassReset)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to retrieve the verification data from db", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to reset password. Please try again later"}, w)
		return
	}

	if verificationData.Code != passResetReq.Code {
		// we should never be here.
		ah.logger.Error("verification code did not match even after verifying PassReset")
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to reset password. Please try again later"}, w)
		return
	}

	if passResetReq.Password != passResetReq.PasswordRe {
		ah.logger.Error("password and password re-enter did not match")
		w.WriteHeader(http.StatusNotAcceptable)
		data.ToJSON(&GenericResponse{Status: false, Message: "Password and re-entered Password are not same"}, w)
		return
	}

	hashedPass, err := ah.hashPassword(passResetReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: UserCreationFailed}, w)
		return
	}

	tokenHash := utils.GenerateRandomString(15)
	err = ah.db.UpdatePassword(context.Background(), userID, hashedPass, tokenHash)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to update user password in db", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Password and re-entered Password are not same"}, w)
		return
	}

	// delete the VerificationData from db
	err = ah.db.DeleteVerificationData(context.Background(), verificationData.Email, verificationData.Type)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to delete the verification data", "error", err))
	}

	w.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Status: false, Message: "Password Reset Successfully"}, w)
}
