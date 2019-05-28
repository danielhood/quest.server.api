package services

import (
	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// TriggerService provides a hook into the game engine for player triggered devices
type TriggerService interface {
	Trigger(*entities.Player, *entities.Device) (string, error)
}

// NewTriggerService creates a new TriggerService
func NewTriggerService(pr repositories.PlayerRepo, dr repositories.DeviceRepo) TriggerService {
	return &triggerService{
		playerRepo: pr,
		deviceRepo: dr,
	}
}

type triggerService struct {
	playerRepo repositories.PlayerRepo
	deviceRepo repositories.DeviceRepo
}

func (s *triggerService) Trigger(p *entities.Player, d *entities.Device) (string, error) {
	// TODO: Process game state for user and return device action code

	// Valid action codes are: NONE, UNKNOWN_PLAYER, NO_QUEST, COMPLETED, ACTIVATE
	return "ACTIVATE", nil
}
