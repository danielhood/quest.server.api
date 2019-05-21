package main

import (
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/handlers"
	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/security"
)

func generateDefaultUsers() {
	userRepo := repositories.NewUserRepo()

	if err := userRepo.Load(); err != nil {
		panic(err)
	}

	if users, err := userRepo.GetAll(); err != nil {
		panic(err)
	} else {
		if len(users) == 0 {
			// Initialize default users if no users currently exist
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
		}
	}
}

func createDefaultRoutes() {
	pingHandler := handlers.NewPing()
	tokenHandler := handlers.NewToken()
	objectHandler := handlers.NewObject()

	auth := security.NewAuthentication()

	http.Handle("/ping", pingHandler)
	http.Handle("/token", tokenHandler)
	http.Handle("/object", AddMiddleware(objectHandler, auth.Authenticate))
}

func main() {
	log.Print("Quest server starting")

	log.Print("Generating default users")
	generateDefaultUsers()

	log.Print("Creating routes")
	createDefaultRoutes()

	log.Print("Listening for connections on port 8080")

	// openssl genrsa -out server.key 2048
	certPath := "server.pem"

	// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
	keyPath := "server.key"

	log.Fatal(http.ListenAndServeTLS(":8080", certPath, keyPath, nil))

	log.Print("Terminating")
}

func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
