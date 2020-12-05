package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"time"

	"github.com/d-vignesh/go-jwt-auth/data"
	"github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/crypto/bcrypt"
)

// Authentication interface lists the methods that our authentication service should implement
type Authentication interface {
	Authenticate(reqUser *data.User, user *data.User) bool
	GenerateAccessToken(user *data.User) (string, error)
	GenerateRefreshToken(user *data.User) (string, error)
	GenerateCustomKey(userID string, password string) string
	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, string, error)
}

// RefreshTokenCustomClaims specifies the claims for refresh token
type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}

// AuthService is the implementation of our Authentication
type AuthService struct {
	logger  hclog.Logger
	configs *utils.Configurations
}

// NewAuthService returns a new instance of the auth service
func NewAuthService(logger hclog.Logger, configs *utils.Configurations) *AuthService {
	return &AuthService{logger, configs}
}

// Authenticate checks the user credentials in request against the db and authenticates the request
func (auth *AuthService) Authenticate(reqUser *data.User, user *data.User) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		auth.logger.Debug("password hashes are not same")
		return false
	}
	return true
}

// GenerateRefreshToken generate a new refresh token for the given user
func (auth *AuthService) GenerateRefreshToken(user *data.User) (string, error) {

	cusKey := auth.GenerateCustomKey(user.ID, user.TokenHash)
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		user.ID,
		cusKey,
		tokenType,
		jwt.StandardClaims{
			Issuer: "bookite.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.configs.RefreshTokenPrivateKeyPath)
	if err != nil {
		auth.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateAccessToken generates a new access token for the given user
func (auth *AuthService) GenerateAccessToken(user *data.User) (string, error) {

	userID := user.ID
	tokenType := "access"

	claims := AccessTokenCustomClaims{
		userID,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(auth.configs.JwtExpiration)).Unix(),
			Issuer:    "bookite.auth.service",
		},
	}

	signBytes, err := ioutil.ReadFile(auth.configs.AccessTokenPrivateKeyPath)
	if err != nil {
		auth.logger.Error("unable to read private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		auth.logger.Error("unable to parse private key", "error", err)
		return "", errors.New("could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (auth *AuthService) GenerateCustomKey(userID string, tokenHash string) string {

	// data := userID + tokenHash
	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userID))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

// ValidateAccessToken parses and validates the given access token
// returns the userId present in the token payload
func (auth *AuthService) ValidateAccessToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(auth.configs.AccessTokenPublicKeyPath)
		if err != nil {
			auth.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			auth.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		auth.logger.Error("unable to parse claims", "error", err)
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", errors.New("invalid token: authentication failed")
	}
	return claims.UserID, nil
}

// ValidateRefreshToken parses and validates the given refresh token
// returns the userId and customkey present in the token payload
func (auth *AuthService) ValidateRefreshToken(tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			auth.logger.Error("Unexpected signing method in auth token")
			return nil, errors.New("Unexpected signing method in auth token")
		}
		verifyBytes, err := ioutil.ReadFile(auth.configs.RefreshTokenPublicKeyPath)
		if err != nil {
			auth.logger.Error("unable to read public key", "error", err)
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			auth.logger.Error("unable to parse public key", "error", err)
			return nil, err
		}

		return verifyKey, nil
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
