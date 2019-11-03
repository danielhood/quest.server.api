package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// Audio holds DeviceService structure
type Audio struct {
}

// NewAudio creates new instance of DeviceService
func NewAudio() *Audio {
	return &Audio{}
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

	case "PUT":
		log.Print("/device/audio:PUT")

		// Device PUT funciton requires device or user level access
		if req.Header.Get("QUEST_AUTH_TYPE") != "device" && req.Header.Get("QUEST_AUTH_TYPE") != "user" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		req.ParseMultipartForm(32 << 20)

		// var deviceType string
		// var responseType string

		for key := range req.Form {
			log.Print("key:", key)
			switch key {
			case "devicetype":
				log.Print("devicetype:", req.FormValue(key))
				// 	deviceType = req.FormValue(key)
			case "responsetype":
				log.Print("requesttype:", req.FormValue(key))
				// 	responseType = req.FormValue(key)
			}
		}

		file, _, err := req.FormFile("filedata")
		defer file.Close()

		if err != nil {
			log.Print("unable to parse file data")
			http.Error(w, "Unable to parse file data", http.StatusInternalServerError)
			return
		}

		fileData := bytes.NewBuffer(nil)
		if _, err := io.Copy(fileData, file); err != nil {
			log.Print("unable to copy file data")
			http.Error(w, "Unable to copy file data", http.StatusInternalServerError)
			return
		}

		log.Print("filedata length: ", fileData.Len())

		w.Write(nil)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Audio) enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
