package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/caseyrwebb/go-jwt-auth/app/data"
	"github.com/caseyrwebb/go-jwt-auth/app/models"
	"github.com/caseyrwebb/go-jwt-auth/app/services"
	"github.com/caseyrwebb/go-jwt-auth/app/utils"
)

// RefreshToken handles refresh token request
func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(UserKey{}).(models.User)
	accessToken, err := ah.authService.GenerateAccessToken(&user)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to generate access token", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		// data.ToJSON(&GenericError{Error: err.Error()}, w)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to generate access token.Please try again later"}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	// data.ToJSON(&TokenResponse{AccessToken: accessToken}, w)
	data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "Successfully generated new access token",
		Data:    &TokenResponse{AccessToken: accessToken},
	}, w)
}

// Greet request greet request
func (ah *AuthHandler) Greet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(UserIDKey{}).(string)
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("hello, " + userID))
	data.ToJSON(&GenericResponse{
		Status:  true,
		Message: "hello," + userID,
	}, w)
}

// GeneratePassResetCode generate a new secret code to reset password.
func (ah *AuthHandler) GeneratePassResetCode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(UserIDKey{}).(string)

	user, err := ah.db.GetUserByID(context.Background(), userID)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to get user to generate secret code for password reset", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to send password reset code. Please try again later"}, w)
		return
	}

	// Send verification mail
	from := "dummy@email.com"
	to := []string{user.Email}
	subject := "Password Reset for go-jwt-auth"
	mailType := services.PassReset
	mailData := &services.MailData{
		Username: user.Username,
		Code:     utils.GenerateRandomString(8),
	}

	mailReq := ah.mailService.NewMail(from, to, subject, mailType, mailData)
	err = ah.mailService.SendMail(mailReq)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to send mail", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to send password reset code. Please try again later"}, w)
		return
	}

	// store the password reset code to db
	verificationData := &models.VerificationData{
		Email:     user.Email,
		Code:      mailData.Code,
		Type:      models.PassReset,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(ah.configs.PassResetCodeExpiration)),
	}

	err = ah.db.StoreVerificationData(context.Background(), verificationData)
	if err != nil {
		ah.logger.Error(fmt.Sprintf("%s %s %d", "unable to store password reset verification data", "error", err))
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericResponse{Status: false, Message: "Unable to send password reset code. Please try again later"}, w)
		return
	}

	ah.logger.Debug("successfully mailed password reset code")
	w.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericResponse{Status: true, Message: "Please check your mail for password reset code"}, w)
}
