package router

import (
	"net/http"

	"github.com/caseyrwebb/go-jwt-auth/app/data"
	"github.com/caseyrwebb/go-jwt-auth/app/handlers"
	"github.com/caseyrwebb/go-jwt-auth/app/services"
	"github.com/caseyrwebb/go-jwt-auth/app/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// create a function to return the routes of app router
func InitRoutes(r *mux.Router, db data.GoDB, logger zap.Logger, configs *utils.Configurations) {
	// validator contains all the methods that are need to validate the user json in request
	validator := data.NewValidation()
	// authService contains all methods that help in authorizing a user request
	authService := services.NewAuthService(logger, configs)

	// mailService contains the utility methods to send an email
	mailService := services.NewSGMailService(logger, configs)

	// UserHandler encapsulates all the services related to user
	uh := handlers.NewAuthHandler(db, logger, configs, validator, authService, mailService)

	// register handlers
	postR := r.Methods(http.MethodPost).Subrouter()

	mailR := r.PathPrefix("/verify").Methods(http.MethodPost).Subrouter()
	mailR.HandleFunc("/mail", uh.VerifyMail)
	mailR.HandleFunc("/password-reset", uh.VerifyPasswordReset)
	mailR.Use(uh.MiddlewareValidateVerificationData)

	postR.HandleFunc("/signup", uh.Signup)
	postR.HandleFunc("/login", uh.Login)
	postR.Use(uh.MiddlewareValidateUser)

	// used the PathPrefix as workaround for scenarios where all the
	// get requests must use the ValidateAccessToken middleware except
	// the /refresh-token request which has to use ValidateRefreshToken middleware
	refToken := r.PathPrefix("/refresh-token").Subrouter()
	refToken.HandleFunc("", uh.RefreshToken)
	refToken.Use(uh.MiddlewareValidateRefreshToken)

	getR := r.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/greet", uh.Greet)
	getR.HandleFunc("/get-password-reset-code", uh.GeneratePassResetCode)
	getR.Use(uh.MiddlewareValidateAccessToken)

	putR := r.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/update-username", uh.UpdateUsername)
	putR.HandleFunc("/reset-password", uh.ResetPassword)
	putR.Use(uh.MiddlewareValidateAccessToken)
}
