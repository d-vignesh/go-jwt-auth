package handlers

import (
	"context"
	"net/http"
	"strings"
	"errors"

	"github.com/d-vignesh/go-jwt-auth/data"
)

// MiddlewareValidateUser validates the user in the request
func (uh *UserHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uh.logger.Debug("user json", r.Body)
		user := &data.User{}

		err := data.FromJSON(user, r.Body)
		if err != nil {
			uh.logger.Error("deserialization of user json failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		// validate the user
		errs := uh.validator.Validate(user)
		if len(errs) != 0 {
			uh.logger.Error("validation of user json failed", "error", errs)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			return
		}

		// add the user to the context
		ctx := context.WithValue(r.Context(), UserKey{}, *user)
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateAccessToken validates whether the request contains a bearer token
// it also decodes and authenticates the given token
func (uh *UserHandler) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		uh.logger.Debug("validating access token")

		token, err := extractToken(r)
		if err != nil {
			uh.logger.Error("Token not provided or malformed")
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		uh.logger.Debug("token present in header", token)
		
		userID, err := uh.authService.ValidateAccessToken(token)
		if err != nil {
			uh.logger.Error("token validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		uh.logger.Debug("access token validated")

		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateAccessToken validates whether the request contains a bearer token
// it also decodes and authenticates the given token
func (uh *UserHandler) MiddlewareValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uh.logger.Debug("validation refresh token")
		uh.logger.Debug("auth header", r.Header.Get("Authorization"))
		token, err := extractToken(r)
		if err != nil {
			uh.logger.Error("token not provided or malformed")
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		uh.logger.Debug("token present in header", token)

		userID, customKey, err := uh.authService.ValidateRefreshToken(token)
		if err != nil {
			uh.logger.Error("token validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}
		uh.logger.Debug("refresh token validated")

		user, err := uh.repo.GetUserByID(context.Background(), userID)
		if err != nil {
			uh.logger.Error("invalid token: wrong userID while parsing", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "invalid token: authentication failed"}, w)
			return
		}

		actualCustomKey := uh.authService.GenerateCustomKey(user.ID, user.Password)
		if customKey != actualCustomKey {
			uh.logger.Debug("wrong token: authetincation failed")
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "invalid toke: authentication failed"}, w)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, *user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("Token not provided or malformed")
	}
	return authHeaderContent[1], nil
}