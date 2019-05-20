package repositories

import (
	"errors"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

// TODO: Move this to somewhere better (redis?)
var players map[uint]*entities.Player

func init() {
	players = make(map[uint]*entities.Player)
}

type PlayerRepo interface {
	GetAll() ([]entities.Player, error)
	Get(id uint) (*entities.Player, error)
	Add(o *entities.Player) error
}

type playerRepo struct {
}

func NewPlayerRepo() PlayerRepo {
	return &playerRepo{}
}

func (r *playerRepo) GetAll() ([]entities.Player, error) {
	allPlayers := make([]entities.Player, len(players))

	idx := 0
	for _, value := range players {
		allPlayers[idx] = *value
		idx++
	}

	return allPlayers, nil
}

func (r *playerRepo) Get(id uint) (*entities.Player, error) {
	if val, ok := players[id]; ok {
		return val, nil
	}

	return nil, errors.New("Player for id not found")
}

func (r *playerRepo) Add(p *entities.Player) error {
	log.Print("Add Player: ", p.Name, " Code: ", p.Code)

	existing, _ := r.Get(p.ID)
	if existing != nil {
		// merge
		existing.Name = p.Name
		existing.QuestID = p.QuestID
		existing.QuestState = p.QuestState
		existing.Achievements = p.Achievements
		existing.Isnabled = p.Isnabled

		return nil
	}

	players[p.ID] = p

	return nil
}
