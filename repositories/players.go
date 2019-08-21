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
	GetByCode(code int) (*entities.Player, error)
	Add(o *entities.Player) error
	Delete(o *entities.Player) error
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

func (r *playerRepo) GetByCode(code int) (*entities.Player, error) {
	for i, p := range players {
		if p.Code == code {
			return &players[i], nil
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
		existing.QuestKey = p.QuestKey
		existing.QuestState = p.QuestState
		existing.QuestStatus = p.QuestStatus
		existing.Achievements = p.Achievements
		existing.IsEnabled = p.IsEnabled
	} else {
		players = append(players, *p)
	}

	return r.store()
}

func (r *playerRepo) Delete(p *entities.Player) error {
	log.Print("Delete Player: ", p.Code)

	for i, player := range players {
		if player.Code == p.Code {
			players[i] = players[len(players)-1]
			players = players[:len(players)-1]
			return r.store()
		}
	}

	return nil
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
	playersJSON, err := r.storageManager.Load("players")

	if err != nil {
		return err
	}

	if len(playersJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(playersJSON), &players); err != nil {
		return err
	}

	log.Printf("Loaded %v player(s)", len(players))

	return nil
}
