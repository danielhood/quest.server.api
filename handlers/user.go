package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/services"
)

// User holds UserService structure
type User struct {
	svc services.UserService
}

// NewUser creates new instance of UserService
func NewUser(ur repositories.UserRepo) *User {
	return &User{services.NewUserService(ur)}
}

func (h *User) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// All user funcitons require user level access
	if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "GET":
		log.Print("/object:GET")

		log.Print("GET params were:", req.URL.Query())

		username := req.URL.Query().Get("username")

		if len(username) == 0 {
			users, _ := h.svc.ReadAll()
			usersBytes, _ := json.Marshal(users)
			w.Write(usersBytes)
		} else {
			user, _ := h.svc.Read(username)
			userBytes, _ := json.Marshal(user)
			w.Write(userBytes)
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
