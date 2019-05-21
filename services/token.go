package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/danielhood/quest.server.api/entities"
)

// Set our secret.
var mySigningKey = []byte("abxx002151!")

// Token defines a token for our application
type Token string

// TokenService provides a token
type TokenService interface {
	GetUserToken(u *entities.User) (string, error)
	GetDeviceToken(d *entities.Device) (string, error)
}

type tokenService struct {
}

type UserClaims struct {
	IsAdmin  bool   `json:"isadmin"`
	AuthType string `json:"authtype"`
	jwt.StandardClaims
}

// NewTokenService creates a new UserService
func NewTokenService() TokenService {
	return &tokenService{}
}

// GetUserToken retrieves a token for a user
func (s *tokenService) GetUserToken(u *entities.User) (string, error) {
	// Set token claims
	claims := UserClaims{
		u.HasRole("AdministratorRole"),
		"user",
		jwt.StandardClaims{
			Id:        strconv.Itoa((int)(u.ID)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "token-service",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	fmt.Printf("Generated User Token for %v %v: %v", u.FirstName, u.LastName, tokenString)

	return tokenString, nil
}

// GetDeviceToken retrieves a token for a device
func (s *tokenService) GetDeviceToken(d *entities.Device) (string, error) {
	// Set token claims
	claims := UserClaims{
		false,
		"device",
		jwt.StandardClaims{
			Id:        strconv.Itoa((int)(d.ID)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "token-service",
		},
	}

	tokenString, err := s.createToken(claims)

	if err == nil {
		fmt.Printf("Generated Device Token for %v %v: %v", d.Hostname, d.Key, tokenString)
	}

	return tokenString, err
}

func (s *tokenService) createToken(claims UserClaims) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	return tokenString, nil
}
