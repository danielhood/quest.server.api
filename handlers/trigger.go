package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/services"
)

// Trigger holds TriggerService structure
type Trigger struct {
	svc services.TriggerService
}

// TriggerRequest holds request parameters for trigger POST.
type TriggerRequest struct {
	PlayerCode string `json:"playercode"`
	DeviceType string `json:"devicetype"`
}

// NewTrigger creates new instance of TriggerService
func NewTrigger(pr repositories.PlayerRepo, dr repositories.DeviceRepo) *Trigger {
	return &Trigger{services.NewTriggerService(pr, dr)}
}

func (h *Trigger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// All trigger funcitons require device level access
	if req.Header.Get("QUEST_AUTH_TYPE") != "device" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "POST":
		log.Print("/trigger:POST")

		var triggerRequest = h.parseRequest(w, req)

		if triggerRequest == nil {
			return
		}

		var deviceActionCode, err = h.svc.Trigger(triggerRequest.PlayerCode, triggerRequest.DeviceType)

		if err != nil {
			http.Error(w, "Unable to process trigger", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(deviceActionCode))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Trigger) parseRequest(w http.ResponseWriter, req *http.Request) *TriggerRequest {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		http.Error(w, "Empty TriggerRequest passed", http.StatusInternalServerError)
		return nil
	}

	var triggerRequest TriggerRequest
	if err = json.Unmarshal(requestBody, &triggerRequest); err != nil {
		http.Error(w, "Unable to parse TriggerRequest json", http.StatusInternalServerError)
		return nil
	}

	if len(triggerRequest.PlayerCode) == 0 {
		http.Error(w, "PlayerCode not specified", http.StatusInternalServerError)
		return nil
	}

	if len(triggerRequest.DeviceType) == 0 {
		http.Error(w, "DeviceType not specified", http.StatusInternalServerError)
		return nil
	}

	return &triggerRequest
}
