package services

import (
	"github.com/danielhood/quest.server.api/quests"
	"github.com/danielhood/quest.server.api/repositories"
)

// TriggerService provides a hook into the game engine for player triggered devices
type TriggerService interface {
	Trigger(int, string) (string, error)
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

func (s *triggerService) Trigger(playerCode int, deviceType string) (string, error) {
	player, _ := s.playerRepo.GetByCode(playerCode)

	if player == nil {
		return quests.QuestResponseUnknownPlayer, nil
	}

	if !player.IsEnabled {
		return quests.QuestResponseUnknownPlayer, nil
	}

	if player.QuestKey == "" {
		return quests.QuestResponseNoQuest, nil
	}

	switch player.QuestKey {
	case "FIND_ALL_TREASURE":
		findAllTreasureQuest := quests.NewFindAllTreasureQuest(s.playerRepo)
		return findAllTreasureQuest.Trigger(player, deviceType)
	case "STARS_ORDERED":
		starsOrderedQuest := quests.NewStarsOrderedQuest(s.playerRepo)
		return starsOrderedQuest.Trigger(player, deviceType)
	default:
		player.QuestState = ""
		player.QuestStatus = ""
		s.playerRepo.Add(player)
		return quests.QuestResponseUnknownQuest, nil
	}
}
