package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/entities"
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
		log.Print("/user:GET")

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

	case "POST":
		log.Print("/user:POST")

		var user = h.parsePutRequest(w, req)

		if user == nil {
			return
		}

		_ = h.svc.Create(user)
		userBytes, _ := json.Marshal(user)
		w.Write(userBytes)

	case "PUT":
		log.Print("/user:PUT")

		var user = h.parsePutRequest(w, req)

		if user == nil {
			return
		}

		_ = h.svc.Update(user)
		userBytes, _ := json.Marshal(user)
		w.Write(userBytes)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *User) parsePutRequest(w http.ResponseWriter, req *http.Request) *entities.User {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		http.Error(w, "Empty User passed", http.StatusInternalServerError)
		return nil
	}

	var user entities.User
	if err = json.Unmarshal(requestBody, &user); err != nil {
		http.Error(w, "Unable to parse User json", http.StatusInternalServerError)
		return nil
	}

	if len(user.Username) == 0 {
		http.Error(w, "Username not specified", http.StatusInternalServerError)
		return nil
	}

	return &user
}
