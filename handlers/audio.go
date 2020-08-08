package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/services"
)

// Audio holds DeviceService structure
type Audio struct {
	svc services.AudioService
}

// AudioGetRequest holds request parameters for audio GET.
type AudioGetRequest struct {
	DeviceType   string `json:"devicetype"`
	ResponseType string `json:"responsetype"`
}

// NewAudio creates new instance of DeviceService
func NewAudio(ar repositories.AudioRepo) *Audio {
	return &Audio{services.NewAudioService(ar)}
}

func (h *Audio) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		log.Print("/device/audio:GET")

		// Audio GET funciton requires device or user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" && req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		audioRequest, _ := h.parseGetRequest(w, req)

		if audioRequest == nil {
			return
		}

		w.Write(nil)

	case "PUT":
		log.Print("/device/audio:PUT")

		// Audio PUT funciton requires device or user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" && req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		audio, fileBuffer := h.parsePutRequest(w, req)

		if audio == nil {
			return
		}

		_ = h.svc.Update(audio, fileBuffer)

		w.Write(nil)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Audio) parseGetRequest(w http.ResponseWriter, req *http.Request) (*AudioGetRequest, *bytes.Buffer) {
	http.Error(w, "GET not implemented", http.StatusInternalServerError)
	return nil, nil
}

func (h *Audio) parsePutRequest(w http.ResponseWriter, req *http.Request) (*entities.Audio, *bytes.Buffer) {
	req.ParseMultipartForm(32 << 20)

	var audio entities.Audio

	for key := range req.Form {
		log.Print("key:", key)
		switch key {
		case "devicetype":
			log.Print("devicetype:", req.FormValue(key))
			audio.DeviceType = req.FormValue(key)
		case "responsetype":
			log.Print("responsetype:", req.FormValue(key))
			audio.ResponseType = req.FormValue(key)
		}
	}

	file, _, err := req.FormFile("filedata")
	defer file.Close()

	if err != nil {
		log.Print("unable to parse file data")
		http.Error(w, "Unable to parse file data", http.StatusInternalServerError)
		return nil, nil
	}

	fileData := bytes.NewBuffer(nil)
	if _, err := io.Copy(fileData, file); err != nil {
		log.Print("unable to copy file data")
		http.Error(w, "Unable to copy file data", http.StatusInternalServerError)
		return nil, nil
	}

	log.Print("filedata length: ", fileData.Len())

	return &audio, fileData
}

func (h *Audio) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
