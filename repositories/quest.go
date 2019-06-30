package repositories

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

// TODO: Move this to somewhere better (redis?)
var quests []entities.Quest

func init() {
	quests = make([]entities.Quest, 0)
}

// QuestRepo defines interface
type QuestRepo interface {
	GetAll() ([]entities.Quest, error)
	GetByKey(questKey string) (*entities.Quest, error)
	Add(o *entities.Quest) error
	Delete(o *entities.Quest) error
}

type questRepo struct {
	storageManager StorageManager
}

// NewQuestRepo returns new instance
func NewQuestRepo(sm StorageManager) QuestRepo {
	r := questRepo{
		storageManager: sm,
	}

	r.load()

	return &r
}

func (r *questRepo) GetAll() ([]entities.Quest, error) {
	allQuests := make([]entities.Quest, len(quests))

	idx := 0
	for _, value := range quests {
		allQuests[idx] = value
		idx++
	}

	return allQuests, nil
}

func (r *questRepo) GetByKey(questKey string) (*entities.Quest, error) {
	for i, p := range quests {
		if p.Key == questKey {
			return &quests[i], nil
		}
	}

	return nil, errors.New("Quest for key not found")
}

func (r *questRepo) Add(p *entities.Quest) error {
	log.Print("Add Quest: ", p.Name, " Key: ", p.Key)

	existing, _ := r.GetByKey(p.Key)
	if existing != nil {
		// merge
		existing.Name = p.Name
		existing.Key = p.Key
		existing.Desc = p.Desc
		existing.IsEnabled = p.IsEnabled
	} else {
		quests = append(quests, *p)
	}

	return r.store()
}

func (r *questRepo) Delete(p *entities.Quest) error {
	log.Print("Delete Quest: ", p.Key)

	for i, quest := range quests {
		if quest.Key == p.Key {
			quests[i] = quests[len(quests)-1]
			quests = quests[:len(quests)-1]
			return r.store()
		}
	}

	return nil
}

// Store saves data to redis
func (r *questRepo) store() error {
	log.Print("Saving quests")

	questsJSON, err := json.Marshal(quests)
	if err != nil {
		return err
	}

	return r.storageManager.Store("quests", questsJSON)
}

// Load retrieves data from redis
func (r *questRepo) load() error {
	questsJSON, err := r.storageManager.Load("quests")

	if err != nil {
		return err
	}

	if len(questsJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(questsJSON), &quests); err != nil {
		return err
	}

	log.Printf("Loaded %v quest(s)", len(quests))

	return nil
}
