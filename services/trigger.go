package services

import (
	"github.com/danielhood/quest.server.api/repositories"
)

// TriggerService provides a hook into the game engine for player triggered devices
type TriggerService interface {
	Trigger(string, string) (string, error)
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

func (s *triggerService) Trigger(playerCode string, deviceType string) (string, error) {
	// Valid action codes are: NONE, UNKNOWN_PLAYER, NO_QUEST, COMPLETED, ACTIVATE

	player, _ := s.playerRepo.GetByCode(playerCode)

	if player == nil {
		return "UNKNOWN_PLAYER", nil
	}

	if !player.IsEnabled {
		return "UNKNOWN_PLAYER", nil
	}

	if player.QuestKey == "" {
		return "NO_QUEST", nil
	}

	// TODO: Process game state for player on current quest and return device action code

	return "ACTIVATE", nil
}
