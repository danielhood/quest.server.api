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

// Player holds PlayerService structure
type Player struct {
	svc services.PlayerService
}

// NewPlayer creates new instance of PlayerService
func NewPlayer(ur repositories.PlayerRepo) *Player {
	return &Player{services.NewPlayerService(ur)}
}

func (h *Player) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// All player funcitons require user level access
	if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "GET":
		log.Print("/player:GET")

		log.Print("GET params were:", req.URL.Query())

		playerCode := req.URL.Query().Get("code")

		if len(playerCode) == 0 {
			players, _ := h.svc.ReadAll()
			playersBytes, _ := json.Marshal(players)
			w.Write(playersBytes)
		} else {
			player, err := h.svc.Read(playerCode)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			playerBytes, _ := json.Marshal(player)
			w.Write(playerBytes)
		}

	case "POST":
		log.Print("/player:POST")

		var player = h.parsePutRequest(w, req)

		if player == nil {
			return
		}

		_ = h.svc.Create(player)
		playerBytes, _ := json.Marshal(player)
		w.Write(playerBytes)

	case "PUT":
		log.Print("/player:PUT")

		var player = h.parsePutRequest(w, req)

		if player == nil {
			return
		}

		_ = h.svc.Update(player)
		playerBytes, _ := json.Marshal(player)
		w.Write(playerBytes)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Player) parsePutRequest(w http.ResponseWriter, req *http.Request) *entities.Player {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		http.Error(w, "Empty Player passed", http.StatusInternalServerError)
		return nil
	}

	var player entities.Player
	if err = json.Unmarshal(requestBody, &player); err != nil {
		http.Error(w, "Unable to parse Player json", http.StatusInternalServerError)
		return nil
	}

	if len(player.Code) == 0 {
		http.Error(w, "code not specified", http.StatusInternalServerError)
		return nil
	}

	return &player
}
