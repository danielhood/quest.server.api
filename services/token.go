package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// Set our secret.
var mySigningKey = []byte("abxx002151!")

// Token defines a token for our application
type Token string

// TokenService provides a token
type TokenService interface {
	ProcessUserLogin(username string, password string) (string, error)
	ProcessDeviceLogin(hostname string, deviceKey string) (string, error)
}

type tokenService struct {
	userRepo   repositories.UserRepo
	deviceRepo repositories.DeviceRepo
}

type userClaims struct {
	IsAdmin  bool   `json:"isadmin"`
	AuthType string `json:"authtype"`
	jwt.StandardClaims
}

// NewTokenService creates a new UserService
func NewTokenService(ur repositories.UserRepo, dr repositories.DeviceRepo) TokenService {
	return &tokenService{
		userRepo:   ur,
		deviceRepo: dr,
	}
}

func (s *tokenService) ProcessUserLogin(username string, password string) (string, error) {
	log.Print("Request User: ", username)

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		log.Print("Error retrieving username: ", err)
		return "", err
	}

	if password != user.Password {
		log.Print("Invalid password")
		return "", errors.New("Invalid password")
	}

	if !user.IsEnabled {
		return "", errors.New("User not enabled")
	}

	return s.getUserToken(user)
}

func (s *tokenService) ProcessDeviceLogin(hostname string, deviceKey string) (string, error) {
	log.Print("Request Hostname: ", hostname, " DeviceKey: ", deviceKey)

	device, err := s.deviceRepo.GetByHostnameAndKey(hostname, deviceKey)

	if err != nil {
		// Device doesn't exist yet
		log.Print("Registering new device")

		newDevice := entities.Device{
			Hostname:     hostname,
			DeviceKey:    deviceKey,
			IsRegistered: false,
			IsEnabled:    false,
		}

		s.deviceRepo.Add(&newDevice)

		return "", err
	}

	// Check if registered
	if !device.IsRegistered {
		// Device exists, but not registered, so return no token
		return "", errors.New("Device not registered")

	}

	return s.getDeviceToken(device)
}

// GetUserToken retrieves a token for a user
func (s *tokenService) getUserToken(u *entities.User) (string, error) {
	// Set token claims
	claims := userClaims{
		u.HasRole("AdministratorRole"),
		"user",
		jwt.StandardClaims{
			Id:        u.Username,
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
func (s *tokenService) getDeviceToken(d *entities.Device) (string, error) {
	// Set token claims
	claims := userClaims{
		false,
		"device",
		jwt.StandardClaims{
			Id:        d.Hostname + ":" + d.DeviceKey,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "token-service",
		},
	}

	tokenString, err := s.createToken(claims)

	if err == nil {
		fmt.Printf("Generated Device Token for %v %v: %v", d.Hostname, d.DeviceKey, tokenString)
	}

	return tokenString, err
}

func (s *tokenService) createToken(claims userClaims) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	return tokenString, nil
}
