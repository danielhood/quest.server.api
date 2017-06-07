// Ref: https://github.com/komand/gosea

package main

import (
  "log"
  "net/http"
  "github.com/danielhood/loco.server/handlers"
  //"github.com/danielhood/loco.server/security"
)

func main () {
  log.Print("Loco server starting...")

  // openssl genrsa -out server.key 2048
  certPath := "server.pem"

  // openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
  keyPath := "server.key"

  pingHandler := handlers.NewPing()
  tokenHandler := handlers.NewToken()
  //testHandler := handlers.NewTest()

  //auth := security.NewAuthentication()

  http.Handle("/ping", pingHandler)
  http.Handle("/token", tokenHandler)
  //http.Handle("/test", AddMiddleware(testHandler, auth.Authenticate))

  log.Print("Startup complete")


  log.Fatal(http.ListenAndServeTLS(":8080", certPath, keyPath, nil))

}

func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
