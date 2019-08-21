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
	HasStarRed    bool `json:"hasstarred"`
	HasStarYellow bool `json:"hasstaryellow"`
	HasStarGreen  bool `json:"hasstargreen"`
	HasStarBlue   bool `json:"hasstarblue"`
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
		if questState.HasStarRed {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.HasStarRed = true
		}
		break
	case DeviceTypeStarYellow:
		if !questState.HasStarRed {
			triggerResponse = QuestResponseItemCollectedTooSoon
		} else if questState.HasStarYellow {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.HasStarYellow = true
		}
		break
	case DeviceTypeStarGreen:
		if !questState.HasStarYellow {
			triggerResponse = QuestResponseItemCollectedTooSoon
		} else if questState.HasStarGreen {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.HasStarGreen = true
		}
		break
	case DeviceTypeStarBlue:
		if !questState.HasStarGreen {
			triggerResponse = QuestResponseItemCollectedTooSoon
		} else if questState.HasStarBlue {
			triggerResponse = QuestResponesItemAlreadyCollected
		} else {
			questState.HasStarBlue = true
			triggerResponse = QuestResponseCompleted
		}
		break
	default:
		triggerResponse = QuestResponseItemNotPartOfQuest
	}

	if questState.HasStarRed && questState.HasStarYellow && questState.HasStarGreen && questState.HasStarBlue {
		player.QuestStatus = QuestStatusCompleted
		log.Print("QuestStatus: ", QuestStatusCompleted)
		player.Achievements = append(player.Achievements, QuestKeyStarsOrdered)
	} else {
		player.QuestStatus = QuestStatusActive
	}

	questStateBytes, _ := json.Marshal(questState)
	player.QuestState = string(questStateBytes)

	q.playerRepo.Add(player)

	return triggerResponse, nil
}
