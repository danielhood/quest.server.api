package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/services"
)

// Token contains strucutre of a token handler
type Token struct {
	Service  services.TokenService
	userRepo repositories.UserRepo
}

// TokenRequest holds request parameters for new token.
// Admin Tokens can be requested by a registered username/password pair.
// Device Tokens can be requested by a registered hostname/devicekey pair.
type TokenRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Hostname  string `json:"hostname"`
	DeviceKey string `json:"devicekey"`
}

// NewToken creates new handler for tokens
func NewToken() *Token {
	return &Token{
		services.NewTokenService(),
		repositories.NewUserRepo(),
	}
}

// Handler will return tokens
func (t *Token) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Print("/token:GET")

		//log.Print("GET params were:", req.URL.Query())

		requestBody, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
			return
		}

		if len(requestBody) == 0 {
			http.Error(w, "Empty TokenRequest passed", http.StatusInternalServerError)
			return
		}

		//log.Print("GET body was:", requestBody, " Length: ", len(requestBody))

		var tokenRequest TokenRequest
		if err = json.Unmarshal(requestBody, &tokenRequest); err != nil {
			http.Error(w, "Unable to parse token request json", http.StatusInternalServerError)
			return
		}

		if len(tokenRequest.Username) > 0 {
			t.processUserLogin(w, tokenRequest)
		} else {
			t.processDeviceLogin(w, tokenRequest)
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (t *Token) processUserLogin(w http.ResponseWriter, tokenRequest TokenRequest) {
	log.Print("Request User: ", tokenRequest.Username)

	user, err := t.userRepo.GetByUsername(tokenRequest.Username)
	if err != nil {
		http.Error(w, "Failed to verify user credentials", http.StatusInternalServerError)
		return
	}

	if tokenRequest.Password != user.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := t.Service.Get(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token))
}

func (t *Token) processDeviceLogin(w http.ResponseWriter, tokenRequest TokenRequest) {
	log.Print("Request Hostname: ", tokenRequest.Hostname, " DeviceKey: ", tokenRequest.DeviceKey)

	http.Error(w, "Device login not supported", http.StatusInternalServerError)
}
