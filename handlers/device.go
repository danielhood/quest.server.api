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

// Device holds DeviceService structure
type Device struct {
	svc services.DeviceService
}

// DeviceGetRequest holds request parameters for device GET.
type DeviceGetRequest struct {
	Hostname  string `json:"hostname"`
	DeviceKey string `json:"devicekey"`
}

// NewDevice creates new instance of DeviceService
func NewDevice(dr repositories.DeviceRepo) *Device {
	return &Device{services.NewDeviceService(dr)}
}

func (h *Device) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.enableCors(&w)

	switch req.Method {
	case "OPTIONS":
		log.Print("/token:OPTIONS")
		if req.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Allow", req.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
		}
		w.Header().Set("Access-Control-Allow-Headers", "authorization,access-control-allow-origin,content-type")

	case "GET":

		// Device GET funciton requires device or user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" && req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/device:GET")
		var deviceHostname string
		var deviceKey string
		deviceHostname = req.URL.Query().Get("hostname")
		deviceKey = req.URL.Query().Get("key")

		if len(deviceHostname) == 0 {
			var deviceGetRequest = h.parseGetRequest(w, req)

			if deviceGetRequest == nil {
				return
			}

			deviceHostname = deviceGetRequest.Hostname
			deviceKey = deviceGetRequest.DeviceKey
		}

		var deviceBytes []byte
		if len(deviceHostname) == 0 {
			// Device GET all funciton requires user level access
			if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			deviceList, _ := h.svc.ReadAll()
			deviceBytes, _ = json.Marshal(deviceList)
		} else {
			device, err := h.svc.Read(deviceHostname, deviceKey)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			deviceBytes, _ = json.Marshal(device)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(deviceBytes)

	case "PUT":

		// Device PUT funciton requires user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Print("/device:PUT")

		var device = h.parsePutRequest(w, req)

		if device == nil {
			return
		}

		_ = h.svc.Update(device)
		deviceBytes, _ := json.Marshal(device)

		w.Header().Set("Content-Type", "application/json")
		w.Write(deviceBytes)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Device) parseGetRequest(w http.ResponseWriter, req *http.Request) *DeviceGetRequest {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		return &DeviceGetRequest{}
	}

	var deviceGetRequest DeviceGetRequest
	if err = json.Unmarshal(requestBody, &deviceGetRequest); err != nil {
		http.Error(w, "Unable to parse DeviceGetRequest json", http.StatusInternalServerError)
		return nil
	}

	if len(deviceGetRequest.Hostname) == 0 {
		http.Error(w, "Hostname not specified", http.StatusInternalServerError)
		return nil
	}

	if len(deviceGetRequest.DeviceKey) == 0 {
		http.Error(w, "DeviceKey not specified", http.StatusInternalServerError)
		return nil
	}

	return &deviceGetRequest
}

func (h *Device) parsePutRequest(w http.ResponseWriter, req *http.Request) *entities.Device {
	requestBody, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusInternalServerError)
		return nil
	}

	if len(requestBody) == 0 {
		http.Error(w, "Empty Device passed", http.StatusInternalServerError)
		return nil
	}

	var device entities.Device
	if err = json.Unmarshal(requestBody, &device); err != nil {
		http.Error(w, "Unable to parse Device json", http.StatusInternalServerError)
		return nil
	}

	if len(device.Hostname) == 0 {
		http.Error(w, "Hostname not specified", http.StatusInternalServerError)
		return nil
	}

	if len(device.DeviceKey) == 0 {
		http.Error(w, "DeviceKey not specified", http.StatusInternalServerError)
		return nil
	}

	return &device
}

func (h *Device) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
