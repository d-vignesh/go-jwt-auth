package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/d-vignesh/go-jwt-auth/data"
)

var (
	JWT_SECRETE_KEY = []byte("highlysecretKey")
	JWT_EXPIRATION = 45
)

type Authentication interface {
	Authenticate(reqUser *data.User, user *data.User) bool
	GenerateAccessToken(user *data.User) (string, error)
	GenerateRefreshToken(user *data.User) (string, error)
	GenerateCustomKey(userID string, password string) string
	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, string, error)
}

type RefreshTokenCustomClaims struct {
	UserID string
	CustomKey string
	KeyType string
	jwt.StandardClaims
}

type AccessTokenCustomClaims struct {
	UserID string
	KeyType string
	jwt.StandardClaims
}
