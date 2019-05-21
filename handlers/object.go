package handlers

import (
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/services"
)

type Object struct {
	svc services.Object
}

func NewObject() *Object {
	return &Object{services.NewObject()}
}

func (h *Object) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		log.Print("/object:GET")

		j, _ := h.svc.GetAll()
		log.Print("Objects:%v", string(j))
		w.Write(j)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
