package quests

import (
	"encoding/json"
	"log"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// FindAllTreasureQuest provides logic for TreasureQuest
type FindAllTreasureQuest interface {
	Trigger(*entities.Player, string) (string, error)
}

type findAllTreasureQuest struct {
	playerRepo repositories.PlayerRepo
}

type findAllTreasureQuestState struct {
	hasTreasure1 bool
	hasTreasure2 bool
	hasTreasure3 bool
	hasTreasure4 bool
}

// NewFindAllTreasureQuest creates a new FindAllTreasureQuest
func NewFindAllTreasureQuest(pr repositories.PlayerRepo) FindAllTreasureQuest {
	return &findAllTreasureQuest{
		playerRepo: pr,
	}
}

func (q *findAllTreasureQuest) Trigger(player *entities.Player, deviceType string) (string, error) {
	// Valid action codes are: NONE, UNKNOWN_PLAYER, NO_QUEST, COMPLETED, ACTIVATE, WRONG_ITEM

	if player.QuestStatus == "" {
		player.QuestStatus = "ACTIVE"
		player.QuestState, _ = json.Marshal(&findAllTreasureQuestState{false, false, false, false})
	}

	var questState findAllTreasureQuestState
	json.Unmarshal(player.QuestState, &questState)

	triggerResponse := "ACTIVATE"

	switch deviceType {
	case "TREASURE:1":
		questState.hasTreasure1 = true
		break
	case "TREASURE:2":
		questState.hasTreasure2 = true
		break
	case "TREASURE:3":
		questState.hasTreasure3 = true
		break
	case "TREASURE:4":
		questState.hasTreasure4 = true
		break
	default:
		triggerResponse = "WRONG_ITEM"
	}

	if questState.hasTreasure1 && questState.hasTreasure2 && questState.hasTreasure3 && questState.hasTreasure4 {
		player.QuestStatus = "COMPLETED"
		log.Print("QuestStatus: COMPLETED")
	}

	player.QuestState, _ = json.Marshal(questState)

	q.playerRepo.Add(player)

	return triggerResponse, nil
}
