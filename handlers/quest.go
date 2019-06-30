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

// Quest holds QuestService structure
type Quest struct {
	svc services.QuestService
}

// NewQuest creates new instance of QuestService
func NewQuest(ur repositories.QuestRepo) *Quest {
	return &Quest{services.NewQuestService(ur)}
}

func (h *Quest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":

		// Quest GET requires device or user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" && req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/quest:GET")

		log.Print("GET params were:", req.URL.Query())

		questCode := req.URL.Query().Get("key")

		if len(questCode) == 0 {
			quests, _ := h.svc.ReadAll()
			questsBytes, _ := json.Marshal(quests)
			w.Write(questsBytes)
		} else {
			quest, err := h.svc.Read(questCode)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			questBytes, _ := json.Marshal(quest)
			w.Write(questBytes)
		}

	case "POST":

		// Quest POST requires user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/quest:POST")

		var quest = h.parsePutRequest(w, req)

		if quest == nil {
			return
		}

		_ = h.svc.Create(quest)
		questBytes, _ := json.Marshal(quest)
		w.Write(questBytes)

	case "PUT":

		// Quest POST requires user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/quest:PUT")

		var quest = h.parsePutRequest(w, req)

		if quest == nil {
			return
		}

		_ = h.svc.Update(quest)
		questBytes, _ := json.Marshal(quest)
		w.Write(questBytes)

	case "DELETE":

		// Quest DELETE requires user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/quest:DELETE")

		var quest = h.parsePutRequest(w, req)

		if quest == nil {
			return
		}

		_ = h.svc.Delete(quest)
		questBytes, _ := json.Marshal(quest)
		w.Write(questBytes)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Quest) parsePutRequest(w http.ResponseWriter, req *http.Request) *entities.Quest {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		http.Error(w, "Empty Quest passed", http.StatusInternalServerError)
		return nil
	}

	var quest entities.Quest
	if err = json.Unmarshal(requestBody, &quest); err != nil {
		http.Error(w, "Unable to parse Quest json", http.StatusInternalServerError)
		return nil
	}

	if len(quest.Key) == 0 {
		http.Error(w, "key not specified", http.StatusInternalServerError)
		return nil
	}

	return &quest
}
