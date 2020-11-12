package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/d-vignesh/go-jwt-auth/data"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
)

type AuthService struct{
	logger hclog.Logger
}

func NewAuthService(logger hclog.Logger) *AuthService {
	return &AuthService{logger}
}

// Authenticate checks the user credentials in request against the db and authenticates the request
func (auth *AuthService) Authenticate(reqUser *data.User, user *data.User) bool {

	if reqUser.Email != user.Email {
		auth.logger.Debug("request email and user email did not match")
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		auth.logger.Debug("password hashes are not same")
		return false
	}
	return true
}

// GenerateRefreshToken generate a new refresh token for the given user
func (auth *AuthService) GenerateRefreshToken(user *data.User) (string, error) {

	userID := user.ID 
	password := user.Password
	cusKey := auth.GenerateCustomKey(userID, password)
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims {
		userID,
		cusKey,
		tokenType,
		jwt.StandardClaims {
			Issuer:	   "bookite.auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRETE_KEY)
}

// GenerateAccessToken generates a new access token for the given user
func (auth *AuthService) GenerateAccessToken(user *data.User) (string, error) {
	
	userID := user.ID
	tokenType := "access"

	claims := AccessTokenCustomClaims {
		userID,
		tokenType,
		jwt.StandardClaims {
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(JWT_EXPIRATION)).Unix(),
			Issuer:    "bookite.auth.service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRETE_KEY)
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and password
func (auth *AuthService) GenerateCustomKey(userID string, password string) string {

	data := userID + password
	h := hmac.New(sha256.New, []byte(JWT_SECRETE_KEY))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// ValidateAccessToken parses and validates the given access token
// returns the userId present in the token payload
func (auth *AuthService) ValidateAccessToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			auth.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		return JWT_SECRETE_KEY, nil
	})

	if err != nil {
		auth.logger.Error("unable to parse claims", "error", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "" , errors.New("invalid token: authentication failed")
	}
	return claims.UserID, nil
}

// ValidateRefreshToken parses and validates the given refresh token
// returns the userId and customkey present in the token payload
func (auth *AuthService) ValidateRefreshToken(tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			auth.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		return JWT_SECRETE_KEY, nil
	})

	if err != nil {
		auth.logger.Error("unable to parse claims", "error", err)
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	auth.logger.Debug("ok", ok)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		auth.logger.Debug("could not extract claims from token")
		return "", "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, claims.CustomKey, nil
}

