package main

import (
	"log"
	"net/http"
	"os"

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
				Username:  "admin",
				Password:  "admin",
				FirstName: "Admin",
				LastName:  "User",
				Roles:     []string{entities.AdministratorRole},
				IsOnline:  false,
			})

			userRepo.Add(&entities.User{
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

func createDefaultRoutes(userRepo repositories.UserRepo, playerRepo repositories.PlayerRepo, deviceRepo repositories.DeviceRepo) {
	pingHandler := handlers.NewPing()
	tokenHandler := handlers.NewToken(userRepo, deviceRepo)
	userHandler := handlers.NewUser(userRepo)
	triggerHandler := handlers.NewTrigger(playerRepo, deviceRepo)

	auth := security.NewAuthentication()

	http.Handle("/ping", pingHandler)
	http.Handle("/token", tokenHandler)
	http.Handle("/user", addMiddleware(userHandler, auth.Authenticate))
	http.Handle("/trigger", addMiddleware(triggerHandler, auth.Authenticate))
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
	playerRepo := repositories.NewPlayerRepo(storageManager)
	deviceRepo := repositories.NewDeviceRepo(storageManager)

	log.Print("Generating default users")
	generateDefaultUsers(userRepo)

	log.Print("Creating routes")
	createDefaultRoutes(userRepo, playerRepo, deviceRepo)

	// openssl genrsa -out server.key 2048
	certPath := "server.pem"

	// openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
	keyPath := "server.key"

	if _, err := os.Stat(keyPath); err == nil {
		log.Print("Listening for connections on HTTPS port 8443")
		log.Fatal(http.ListenAndServeTLS(":8443", certPath, keyPath, nil))
	} else if os.IsNotExist(err) {
		log.Print("Listening for connections on HTTP port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	log.Print("Terminating")
}
