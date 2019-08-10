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
	PlayerCode int    `json:"playercode"`
	DeviceType string `json:"devicetype"`
}

// NewTrigger creates new instance of TriggerService
func NewTrigger(pr repositories.PlayerRepo, dr repositories.DeviceRepo) *Trigger {
	return &Trigger{services.NewTriggerService(pr, dr)}
}

func (h *Trigger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.enableCors(&w)

	switch req.Method {
	case "OPTIONS":
		log.Print("/token:OPTIONS")
		if req.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Allow", req.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
		}
		w.Header().Set("Access-Control-Allow-Headers", "authorization,access-control-allow-origin,content-type")

	case "POST":
		log.Print("/trigger:POST")

		// All trigger funcitons require device level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var triggerRequest = h.parseRequest(w, req)

		if triggerRequest == nil {
			log.Print("nil triggerRequest")
			return
		}

		var deviceActionCode, err = h.svc.Trigger(triggerRequest.PlayerCode, triggerRequest.DeviceType)

		if err != nil {
			log.Print("Unable to process trigger")
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

	if len(triggerRequest.DeviceType) == 0 {
		http.Error(w, "DeviceType not specified", http.StatusInternalServerError)
		return nil
	}

	return &triggerRequest
}

func (h *Trigger) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
