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
	svc      services.TokenService
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
func NewToken(ur repositories.UserRepo, dr repositories.DeviceRepo) *Token {
	return &Token{
		svc:      services.NewTokenService(ur, dr),
		userRepo: ur,
	}
}

// Handler will return tokens
func (t *Token) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Print("/token:GET")

		requestBody, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, "Unable to parse request body", http.StatusUnauthorized)
			return
		}

		if len(requestBody) == 0 {
			http.Error(w, "Empty TokenRequest passed", http.StatusUnauthorized)
			return
		}

		var tokenRequest TokenRequest
		if err = json.Unmarshal(requestBody, &tokenRequest); err != nil {
			http.Error(w, "Unable to parse token request json", http.StatusUnauthorized)
			return
		}

		var token string
		if len(tokenRequest.Username) > 0 {
			token, err = t.svc.ProcessUserLogin(tokenRequest.Username, tokenRequest.Password)
		} else if len(tokenRequest.Hostname) > 0 && len(tokenRequest.DeviceKey) > 0 {
			token, err = t.svc.ProcessDeviceLogin(tokenRequest.Hostname, tokenRequest.DeviceKey)
		} else {
			http.Error(w, "Invalid request body: missing requierd keys", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Failed to verify user credentials", http.StatusUnauthorized)
			return
		}

		w.Write([]byte(token))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
