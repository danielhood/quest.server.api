package services

import (
	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// UserService provides a CRUD interface for Users
type UserService interface {
	Create(*entities.User) error
	Read(string) (*entities.User, error)
	Update(*entities.User) error
	Delete(*entities.User) error
}

// NewUserService creates a new UserService
func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepo(),
	}
}

type userService struct {
	userRepo repositories.UserRepo
}

func (s *userService) Create(u *entities.User) error {
	return s.userRepo.Add(u)
}

func (s *userService) Read(username string) (*entities.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) Update(u *entities.User) error {
	return s.userRepo.Add(u)
}

func (s *userService) Delete(u *entities.User) error {
	return s.userRepo.Delete(u)
}
