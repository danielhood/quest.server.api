package main

import (
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/handlers"
	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/security"
)

func main() {
	log.Print("Loco server starting...")

	log.Print("Loading users...")

	userRepo := repositories.NewUserRepo()

	if err := userRepo.Load(); err != nil {
		panic(err)
	}

	userRepo.Add(&entities.User{
		ID:        1,
		Username:  "admin",
		Password:  "admin",
		FirstName: "Admin",
		LastName:  "User",
		Roles:     []string{entities.AdministratorRole},
		IsOnline:  false,
	})

	userRepo.Add(&entities.User{
		ID:        2,
		Username:  "test",
		Password:  "test",
		FirstName: "Test",
		LastName:  "User",
		Roles:     []string{},
		IsOnline:  false,
	})

	log.Print("Checking Users...")
	userRepo.GetAll()

	log.Print("Generating objects...")

	objectRepo := repositories.NewObjectRepo()

	objectRepo.Add(&entities.Object{
		ID:       1,
		Type:     entities.ObjectType1,
		Name:     "Object 1",
		X:        0,
		Y:        0,
		Bearing:  0,
		Velocity: 0,
	})

	objectRepo.Add(&entities.Object{
		ID:       2,
		Type:     entities.ObjectType2,
		Name:     "Object 2",
		X:        0,
		Y:        0,
		Bearing:  0.1,
		Velocity: 0.5,
	})

	log.Print("Creating routes...")

	pingHandler := handlers.NewPing()
	tokenHandler := handlers.NewToken()
	objectHandler := handlers.NewObject()

	auth := security.NewAuthentication()

	http.Handle("/ping", pingHandler)
	http.Handle("/token", tokenHandler)
	http.Handle("/object", AddMiddleware(objectHandler, auth.Authenticate))

	log.Print("Startup complete")

	// openssl genrsa -out server.key 2048
	certPath := "server.pem"

	// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
	keyPath := "server.key"

	log.Fatal(http.ListenAndServeTLS(":8080", certPath, keyPath, nil))

}

func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
