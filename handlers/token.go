package handlers

import (
	"net/http"
	"log"
	"strconv"

  "github.com/danielhood/quest.server.api/services"
	"github.com/danielhood/quest.server.api/repositories"
)

type Token struct {
	Service services.TokenService
	userRepo repositories.UserRepo
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
		log.Print("/token:GETxx")

		log.Print("GET params were:", req.URL.Query())

		userIdStr := req.URL.Query().Get("user")

		if userIdStr == "" {
			http.Error(w, "User not specified", http.StatusInternalServerError)
		}

		// TODO: Lookup user based on login information passed into token get

		userId, err := strconv.ParseUint(userIdStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid User specified", http.StatusInternalServerError)
		}

		user, err := t.userRepo.Get(uint(userId))
		if (err != nil) {
			http.Error(w, "Failed to verify user credentials", http.StatusInternalServerError)
		}

		token, err := t.Service.Get(user)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		}
		w.Write([]byte(token))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}
