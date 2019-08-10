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
	log.Print("Generating default users")
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
				IsEnabled: true,
			})

			userRepo.Add(&entities.User{
				Username:  "test",
				Password:  "test",
				FirstName: "Test",
				LastName:  "User",
				Roles:     []string{},
				IsOnline:  false,
				IsEnabled: false,
			})
		}
	}
}

func generateDefaultQuests(questRepo repositories.QuestRepo) {
	log.Print("Generating default quests")
	if quests, err := questRepo.GetAll(); err != nil {
		panic(err)
	} else {
		if len(quests) == 0 {
			questRepo.Add(&entities.Quest{
				Key:       "FIND_ALL_TREASURE",
				Name:      "Find All Treasure",
				Desc:      "Find all treasure in the kingdom",
				IsEnabled: true,
			})
		}
	}
}

func generateDefaultPlayers(playerRepo repositories.PlayerRepo) {
	log.Print("Generating default players")
	if players, err := playerRepo.GetAll(); err != nil {
		panic(err)
	} else {
		if len(players) == 0 {
			playerRepo.Add(&entities.Player{
				Code:      12345678,
				Name:      "Test User",
				IsEnabled: true,
				QuestKey:  "FIND_ALL_TREASURE",
			})
		}
	}
}

func createDefaultRoutes(userRepo repositories.UserRepo, playerRepo repositories.PlayerRepo, deviceRepo repositories.DeviceRepo, questRepo repositories.QuestRepo) {
	pingHandler := handlers.NewPing()
	tokenHandler := handlers.NewToken(userRepo, deviceRepo)
	userHandler := handlers.NewUser(userRepo)
	playerHandler := handlers.NewPlayer(playerRepo)
	triggerHandler := handlers.NewTrigger(playerRepo, deviceRepo)
	deviceHandler := handlers.NewDevice(deviceRepo)
	questHandler := handlers.NewQuest(questRepo)

	auth := security.NewAuthentication()

	http.Handle("/ping", pingHandler)
	http.Handle("/token", tokenHandler)
	http.Handle("/user", addMiddleware(userHandler, auth.Authenticate))
	http.Handle("/player", addMiddleware(playerHandler, auth.Authenticate))
	http.Handle("/trigger", addMiddleware(triggerHandler, auth.Authenticate))
	http.Handle("/device", addMiddleware(deviceHandler, auth.Authenticate))
	http.Handle("/quest", addMiddleware(questHandler, auth.Authenticate))
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
	questRepo := repositories.NewQuestRepo(storageManager)

	generateDefaultUsers(userRepo)
	generateDefaultQuests(questRepo)
	generateDefaultPlayers(playerRepo)

	log.Print("Creating routes")
	createDefaultRoutes(userRepo, playerRepo, deviceRepo, questRepo)

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
