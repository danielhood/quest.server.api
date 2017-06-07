package handlers

import (
	"net/http"
	"log"

  "github.com/danielhood/loco.server/services"
)

type Token struct {
	Service services.TokenService
}

// NewToken creates new handler for tokens
func NewToken() *Token {
	return &Token{services.NewTokenService()}
}

// Handler will return tokens
func (t *Token) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Print("/token:GET")

		// TODO: Take in login information
		user := &services.User{
			Id:        1,
			FirstName: "Admin",
			LastName:  "User",
			Roles:     []string{services.AdministratorRole},
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
