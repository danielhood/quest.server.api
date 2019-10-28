package services

import (
	"log"

	"github.com/danielhood/quest.server.api/quests"
	"github.com/danielhood/quest.server.api/repositories"
)

// TriggerService provides a hook into the game engine for player triggered devices
type TriggerService interface {
	Trigger(int, string) (string, error)
	GetLastPlayerCode() (int, error)
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

func (s *triggerService) GetLastPlayerCode() (int, error) {
	return s.playerRepo.GetLastPlayerCode()
}

func (s *triggerService) Trigger(playerCode int, deviceType string) (string, error) {
	log.Print("Trigger: playerCode: ", playerCode)

	if deviceType == quests.DeviceTypeAdminTrigger {
		s.playerRepo.SetLastPlayerCode(playerCode)
		return quests.QuestResponseActivate, nil
	}

	player, _ := s.playerRepo.GetByCode(playerCode)

	if player == nil {
		log.Print("ERROR: Unknown player")
		return quests.QuestResponseUnknownPlayer, nil
	}

	log.Print("Found player.code: ", player.Code)

	if !player.IsEnabled {
		return quests.QuestResponseUnknownPlayer, nil
	}

	if player.QuestKey == "" {
		return quests.QuestResponseNoQuest, nil
	}

	switch player.QuestKey {
	case quests.QuestKeyFindAllTreasure:
		findAllTreasureQuest := quests.NewFindAllTreasureQuest(s.playerRepo)
		return findAllTreasureQuest.Trigger(player, deviceType)
	case quests.QuestKeyStarsOrdered:
		starsOrderedQuest := quests.NewStarsOrderedQuest(s.playerRepo)
		return starsOrderedQuest.Trigger(player, deviceType)
	default:
		player.QuestState = ""
		player.QuestStatus = ""
		s.playerRepo.Add(player)
		return quests.QuestResponseUnknownQuest, nil
	}
}
