package quests

import (
	"encoding/json"
	"log"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// StarsOrderedQuest provides logic for TreasureQuest
type StarsOrderedQuest interface {
	Trigger(*entities.Player, string) (string, error)
}

type starsOrderedQuest struct {
	playerRepo repositories.PlayerRepo
}

type starsOrderedQuestState struct {
	hasStarRed    bool
	hasStarYellow bool
	hasStarGreen  bool
	hasStarBlue   bool
}

// NewStarsOrderedQuest creates a new StarsOrderedQuest
func NewStarsOrderedQuest(pr repositories.PlayerRepo) StarsOrderedQuest {
	return &starsOrderedQuest{
		playerRepo: pr,
	}
}

func (q *starsOrderedQuest) hasCompletedQuest(player *entities.Player) bool {
	for _, key := range player.Achievements {
		if key == QuestKeyStarsOrdered {
			return true
		}
	}
	return false
}

func (q *starsOrderedQuest) Trigger(player *entities.Player, deviceType string) (string, error) {
	if q.hasCompletedQuest(player) {
		return QuestResponseCompleted, nil
	}

	if player.QuestStatus == "" {
		player.QuestStatus = QuestStatusActive

		questStateBytes, _ := json.Marshal(&starsOrderedQuestState{false, false, false, false})
		player.QuestState = string(questStateBytes)
	}

	var questState starsOrderedQuestState
	json.Unmarshal([]byte(player.QuestState), &questState)

	triggerResponse := QuestResponseActivate

	switch deviceType {
	case DeviceTypeStarRed:
		if questState.hasStarRed {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.hasStarRed = true
		}
		break
	case DeviceTypeStarYellow:
		if !questState.hasStarRed {
			triggerResponse = QuestResponseItemNotPartOfQuest
		} else if questState.hasStarYellow {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.hasStarYellow = true
		}
		break
	case DeviceTypeStarGreen:
		if !questState.hasStarYellow {
			triggerResponse = QuestResponseItemNotPartOfQuest
		} else if questState.hasStarGreen {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.hasStarGreen = true
		}
		break
	case DeviceTypeStarBlue:
		if !questState.hasStarGreen {
			triggerResponse = QuestResponseItemNotPartOfQuest
		} else if questState.hasStarBlue {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.hasStarBlue = true
		}
		break
	default:
		triggerResponse = QuestResponseItemNotPartOfQuest
	}

	if questState.hasStarRed && questState.hasStarYellow && questState.hasStarGreen && questState.hasStarBlue {
		player.QuestStatus = QuestStatusCompleted
		log.Print("QuestStatus: ", QuestStatusCompleted)
		player.Achievements = append(player.Achievements, QuestKeyStarsOrdered)
	}

	questStateBytes, _ := json.Marshal(questState)
	player.QuestState = string(questStateBytes)

	q.playerRepo.Add(player)

	return triggerResponse, nil
}
