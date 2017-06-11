package repositories

import (
  "errors"
  "log"

  "github.com/danielhood/loco.server/entities"
)

// TODO: Move this to somewhere better (redis?)
var users map[uint]*entities.User

func init() {
  users = make(map[uint]*entities.User)
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

func (r *userRepo)GetAll() ([]entities.User, error) {
  allUsers := make([]entities.User, len(users))

  idx := 0
  for _, value := range users {
    allUsers[idx] = *value
    idx++
  }

  return allUsers, nil
}

func (r *userRepo)Get(id uint) (*entities.User, error) {
  if val, ok := users[id]; ok {
    return val, nil;
  }

  return nil, errors.New("User for id not found")
}

func (r *userRepo)Add(u *entities.User) error {
  log.Print("Add User: ", u.Username)

  existing, _ := r.Get(u.Id)
  if existing != nil {
    // merge only online status for now
    existing.IsOnline = u.IsOnline
    return nil
  }

  users[u.Id] = u

  return nil
}
