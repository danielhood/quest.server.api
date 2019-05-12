package services

import (
  "github.com/danielhood/quest.server.api/entities"
  "github.com/danielhood/quest.server.api/repositories"
)

// UserService provides a CRUD interface for Users
type UserService interface {
	Create(*entities.User) error
	Read(uint) (*entities.User, error)
	Update(*entities.User) error
}

// NewUserService creates a new UserService
func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepo(),
	}
}

type userService struct{
	userRepo repositories.UserRepo
}

func (s *userService) Create(u *entities.User) error {
	return s.userRepo.Add(u)
}

func (s *userService) Read(id uint) (*entities.User, error) {
	return s.userRepo.Get(id)
}

func (s *userService) Update(u *entities.User) error {
	return s.userRepo.Add(u)
}

// TODO: Delete
