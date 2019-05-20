package repositories

// https://github.com/go-redis/redis

import (
	"errors"
	"fmt"
	"log"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/go-redis/redis"
)

// TODO: Move this to somewhere better (redis?)
var users map[uint]*entities.User

var redisClient *redis.Client

func init() {
	users = make(map[uint]*entities.User)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type UserRepo interface {
	GetAll() ([]entities.User, error)
	Get(id uint) (*entities.User, error)
	Add(o *entities.User) error
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (r *userRepo) GetAll() ([]entities.User, error) {
	allUsers := make([]entities.User, len(users))

	userJSON, err := redisClient.Get("users").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("usersJson", userJSON)

	idx := 0
	for _, value := range users {
		allUsers[idx] = *value
		idx++
	}

	return allUsers, nil
}

func (r *userRepo) Get(id uint) (*entities.User, error) {
	if val, ok := users[id]; ok {
		return val, nil
	}

	return nil, errors.New("User for id not found")
}

func (r *userRepo) Add(u *entities.User) error {
	log.Print("Add User: ", u.Username)

	existing, _ := r.Get(u.ID)
	if existing != nil {
		// merge only online status for now
		existing.IsOnline = u.IsOnline
		return nil
	}

	users[u.ID] = u

	return nil
}
