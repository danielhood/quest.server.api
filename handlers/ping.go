package handlers

import (
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/services"
)

// Ping holds handler structure
type Ping struct {
	svc services.Ping
}

// NewPing creates an instance of Ping
func NewPing() *Ping {
	return &Ping{services.NewPing()}
}

func (h *Ping) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Print("/ping:GET")

		s := h.svc.Get()
		w.Write([]byte(s))

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
