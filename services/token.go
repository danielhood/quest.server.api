package services

import (
	"errors"
	"time"
  "fmt"
  "strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/danielhood/loco.server/entities"
)

// Set our secret.
var mySigningKey = []byte("abxx002151!")

// Token defines a token for our application
type Token string

// TokenService provides a token
type TokenService interface {
	Get(u *entities.User) (string, error)
}

type tokenService struct {
}

type UserClaims struct{
  Admin bool  `json:"admin"`
  User  entities.User  `json:"user"`
  jwt.StandardClaims
}

// NewTokenService creates a new UserService
func NewTokenService() TokenService {
	return &tokenService{}
}

// Get retrieves a token for a user
// TODO: Take user credentials and verify them against what's in database
func (s *tokenService) Get(u *entities.User) (string, error) {
	// Set token claims
  claims := UserClaims {
    u.HasRole("AdministratorRole"),
    *u,
    jwt.StandardClaims {
      Id: strconv.Itoa((int)(u.Id)),
      ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
      Issuer: "token-service",
    },
  }

  // Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

  fmt.Printf("Generated Token for %v %v: %v", u.FirstName, u.LastName, tokenString)

	return tokenString, nil
}
