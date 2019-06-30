package services

import (
	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// QuestService provides a CRUD interface for Quests
type QuestService interface {
	Create(*entities.Quest) error
	ReadAll() ([]entities.Quest, error)
	Read(string) (*entities.Quest, error)
	Update(*entities.Quest) error
	Delete(*entities.Quest) error
}

// NewQuestService creates a new QuestService
func NewQuestService(ur repositories.QuestRepo) QuestService {
	return &questService{
		questRepo: ur,
	}
}

type questService struct {
	questRepo repositories.QuestRepo
}

func (s *questService) Create(p *entities.Quest) error {
	return s.questRepo.Add(p)
}

func (s *questService) ReadAll() ([]entities.Quest, error) {
	return s.questRepo.GetAll()
}

func (s *questService) Read(questKey string) (*entities.Quest, error) {
	return s.questRepo.GetByKey(questKey)
}

func (s *questService) Update(p *entities.Quest) error {
	return s.questRepo.Add(p)
}

func (s *questService) Delete(p *entities.Quest) error {
	return s.questRepo.Delete(p)
}
