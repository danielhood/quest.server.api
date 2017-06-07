package handlers

import (
  "net/http"

  "github.com/danielhood/loco.server/services"
)

type Ping struct {
  svc services.Ping
}

func NewPing() *Ping {
  return &Ping{services.NewPing()}
}

func (h *Ping) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case "GET":
        s := h.svc.Get()
        w.Write([]byte(s))
    default:
        http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    }
}
