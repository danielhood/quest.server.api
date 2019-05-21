package repositories

// https://github.com/go-redis/redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/go-redis/redis"
)

var users []entities.User

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Default strucutre
	users = make([]entities.User, 0)
}

type UserRepo interface {
	GetAll() ([]entities.User, error)
	Get(ID uint) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	Add(o *entities.User) error
	Delete(o *entities.User) error
	Load() error
	Store() error
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (r *userRepo) GetAll() ([]entities.User, error) {
	allUsers := make([]entities.User, len(users))

	idx := 0
	for _, value := range users {
		allUsers[idx] = value
		idx++
	}

	return allUsers, nil
}

func (r *userRepo) Get(ID uint) (*entities.User, error) {
	for _, u := range users {
		if u.ID == ID {
			return &u, nil
		}
	}

	return nil, errors.New("User for id not found")
}

func (r *userRepo) GetByUsername(username string) (*entities.User, error) {
	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, errors.New("User for username not found")
}

func (r *userRepo) Add(u *entities.User) error {
	log.Print("Add User: ", u.Username)

	existing, _ := r.Get(u.ID)
	if existing != nil {
		// merge only online status for now
		existing.IsOnline = u.IsOnline
	} else {
		users = append(users, *u)
	}

	return r.Store()
}

func (r *userRepo) Delete(u *entities.User) error {
	log.Print("Delete User: ", u.Username)

	for i, user := range users {
		if user.ID == u.ID {
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			return r.Store()
		}
	}

	return nil
}

// Store saves data to redis
func (r *userRepo) Store() error {
	log.Print("Saving users")

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return err
	}

	log.Print("usersJSON: ", string(usersJSON))

	err = redisClient.Set("users", usersJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Load retrieves data from redis
func (r *userRepo) Load() error {
	log.Print("Loading users")

	userJSON, err := redisClient.Get("users").Result()

	if err != nil {
		return err
	}
	fmt.Println("usersJson", userJSON)

	if len(userJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(userJSON), &users); err != nil {
		return err
	}

	return nil
}
