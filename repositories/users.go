package repositories

// https://github.com/go-redis/redis

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

var users []entities.User

func init() {
	// Default strucutre
	users = make([]entities.User, 0)
}

// UserRepo defines UserRepo interface
type UserRepo interface {
	GetAll() ([]entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	Add(o *entities.User) error
	Delete(o *entities.User) error
}

type userRepo struct {
	storageManager StorageManager
}

// NewUserRepo returns a new UserRepo instance
func NewUserRepo(sm StorageManager) UserRepo {
	ur := userRepo{
		storageManager: sm,
	}

	ur.load()

	return &ur
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

	existing, _ := r.GetByUsername(u.Username)
	if existing != nil {
		// TODO: pull out password into a separate strucutre, and hash it
		existing.IsOnline = u.IsOnline
		existing.FirstName = u.FirstName
		existing.IsEnabled = u.IsEnabled
		existing.LastName = u.LastName
		existing.Password = u.Password
		existing.Roles = u.Roles
	} else {
		users = append(users, *u)
	}

	return r.store()
}

func (r *userRepo) Delete(u *entities.User) error {
	log.Print("Delete User: ", u.Username)

	for i, user := range users {
		if user.Username == u.Username {
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			return r.store()
		}
	}

	return nil
}

func (r *userRepo) store() error {
	log.Print("Saving users")

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return r.storageManager.Store("users", usersJSON)
}

// Load retrieves data from redis
func (r *userRepo) load() error {
	log.Print("Loading users")

	userJSON, err := r.storageManager.Load("users")

	if err != nil {
		return err
	}

	log.Print("usersJson", userJSON)

	if len(userJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(userJSON), &users); err != nil {
		return err
	}

	return nil
}
