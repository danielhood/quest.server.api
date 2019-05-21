package main

import (
	"log"
	"net/http"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/handlers"
	"github.com/danielhood/quest.server.api/repositories"
	"github.com/danielhood/quest.server.api/security"
	"github.com/go-redis/redis"
)

func generateDefaultUsers(userRepo repositories.UserRepo) {
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

func createDefaultRoutes(userRepo repositories.UserRepo) {
	pingHandler := handlers.NewPing()
	tokenHandler := handlers.NewToken(userRepo)
	userHandler := handlers.NewUser(userRepo)

	auth := security.NewAuthentication()

	http.Handle("/ping", pingHandler)
	http.Handle("/token", tokenHandler)
	http.Handle("/user", addMiddleware(userHandler, auth.Authenticate))
}

func addMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

func main() {
	log.Print("Quest server starting")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	storageManager := repositories.NewStorageManager(redisClient)

	userRepo := repositories.NewUserRepo(storageManager)

	log.Print("Generating default users")
	generateDefaultUsers(userRepo)

	log.Print("Creating routes")
	createDefaultRoutes(userRepo)

	log.Print("Listening for connections on port 8080")

	// openssl genrsa -out server.key 2048
	certPath := "server.pem"

	// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
	keyPath := "server.key"

	log.Fatal(http.ListenAndServeTLS(":8080", certPath, keyPath, nil))

	log.Print("Terminating")
}
