package services

import (
	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// PlayerService provides a CRUD interface for Players
type PlayerService interface {
	Create(*entities.Player) error
	ReadAll() ([]entities.Player, error)
	Read(string) (*entities.Player, error)
	Update(*entities.Player) error
	Delete(*entities.Player) error
}

// NewPlayerService creates a new PlayerService
func NewPlayerService(ur repositories.PlayerRepo) PlayerService {
	return &playerService{
		playerRepo: ur,
	}
}

type playerService struct {
	playerRepo repositories.PlayerRepo
}

func (s *playerService) Create(p *entities.Player) error {
	return s.playerRepo.Add(p)
}

func (s *playerService) ReadAll() ([]entities.Player, error) {
	return s.playerRepo.GetAll()
}

func (s *playerService) Read(playerCode string) (*entities.Player, error) {
	return s.playerRepo.GetByCode(playerCode)
}

func (s *playerService) Update(p *entities.Player) error {
	return s.playerRepo.Add(p)
}

func (s *playerService) Delete(p *entities.Player) error {
	return s.playerRepo.Delete(p)
}
