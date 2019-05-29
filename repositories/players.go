package repositories

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

// TODO: Move this to somewhere better (redis?)
var players []entities.Player

func init() {
	players = make([]entities.Player, 0)
}

// PlayerRepo defines interface
type PlayerRepo interface {
	GetAll() ([]entities.Player, error)
	GetByCode(code string) (*entities.Player, error)
	Add(o *entities.Player) error
}

type playerRepo struct {
	storageManager StorageManager
}

// NewPlayerRepo returns new instance
func NewPlayerRepo(sm StorageManager) PlayerRepo {
	r := playerRepo{
		storageManager: sm,
	}

	r.load()

	return &r
}

func (r *playerRepo) GetAll() ([]entities.Player, error) {
	allPlayers := make([]entities.Player, len(players))

	idx := 0
	for _, value := range players {
		allPlayers[idx] = value
		idx++
	}

	return allPlayers, nil
}

func (r *playerRepo) GetByCode(code string) (*entities.Player, error) {
	for _, p := range players {
		if p.Code == code {
			return &p, nil
		}
	}

	return nil, errors.New("Player for id not found")
}

func (r *playerRepo) Add(p *entities.Player) error {
	log.Print("Add Player: ", p.Name, " Code: ", p.Code)

	existing, _ := r.GetByCode(p.Code)
	if existing != nil {
		// merge
		existing.Name = p.Name
		existing.QuestID = p.QuestID
		existing.QuestState = p.QuestState
		existing.Achievements = p.Achievements
		existing.Isnabled = p.Isnabled
	} else {
		players = append(players, *p)
	}

	return r.store()
}

// Store saves data to redis
func (r *playerRepo) store() error {
	log.Print("Saving players")

	playersJSON, err := json.Marshal(players)
	if err != nil {
		return err
	}

	return r.storageManager.Store("players", playersJSON)
}

// Load retrieves data from redis
func (r *playerRepo) load() error {
	log.Print("Loading players")

	playersJSON, err := r.storageManager.Load("players")

	if err != nil {
		return err
	}

	log.Print("playersJSON", playersJSON)

	if len(playersJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(playersJSON), &players); err != nil {
		return err
	}

	return nil
}
