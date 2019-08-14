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

type FindAllTreasureQuestState struct {
	HasTreasure1 bool `json:"hastreasure1"`
	HasTreasure2 bool `json:"hastreasure2"`
	HasTreasure3 bool `json:"hastreasure3"`
	HasTreasure4 bool `json:"hastreasure4"`
}

// NewFindAllTreasureQuest creates a new FindAllTreasureQuest
func NewFindAllTreasureQuest(pr repositories.PlayerRepo) FindAllTreasureQuest {
	return &findAllTreasureQuest{
		playerRepo: pr,
	}
}

func (q *findAllTreasureQuest) hasCompletedQuest(player *entities.Player) bool {
	for _, key := range player.Achievements {
		if key == QuestKeyFindAllTreasure {
			return true
		}
	}
	return false
}

func (q *findAllTreasureQuest) Trigger(player *entities.Player, deviceType string) (string, error) {

	log.Print("In: ", player.QuestState)

	if q.hasCompletedQuest(player) {
		return QuestResponseCompleted, nil
	}

	if player.QuestStatus == "" {
		player.QuestStatus = QuestStatusActive
		log.Print("Resetting quest state")
		questStateBytes, _ := json.Marshal(&FindAllTreasureQuestState{false, false, false, false})
		player.QuestState = string(questStateBytes)
	}

	var questState FindAllTreasureQuestState
	json.Unmarshal([]byte(player.QuestState), &questState)

	log.Print(questState.HasTreasure1, questState.HasTreasure2, questState.HasTreasure3, questState.HasTreasure4)

	triggerResponse := QuestResponseActivate

	switch deviceType {
	case DeviceTypeTreasure1:
		if questState.HasTreasure1 {
			triggerResponse = QuestResponesItemAlreadyCollected
		}
		questState.HasTreasure1 = true
		break
	case DeviceTypeTreasure2:
		if questState.HasTreasure2 {
			triggerResponse = QuestResponesItemAlreadyCollected
		}
		questState.HasTreasure2 = true
		break
	case DeviceTypeTreasure3:
		if questState.HasTreasure3 {
			triggerResponse = QuestResponesItemAlreadyCollected
		}
		questState.HasTreasure3 = true
		break
	case DeviceTypeTreasure4:
		if questState.HasTreasure4 {
			triggerResponse = QuestResponesItemAlreadyCollected
		}
		questState.HasTreasure4 = true
		break
	default:
		triggerResponse = QuestResponseItemNotPartOfQuest
	}

	if questState.HasTreasure1 && questState.HasTreasure2 && questState.HasTreasure3 && questState.HasTreasure4 {
		player.QuestStatus = QuestStatusCompleted
		log.Print("QuestStatus: ", QuestStatusCompleted)
		player.Achievements = append(player.Achievements, QuestKeyFindAllTreasure)
		triggerResponse = QuestResponseCompleted
	}

	questStateBytes, err := json.Marshal(questState)

	if err != nil {
		log.Print(err)
	}

	player.QuestState = string(questStateBytes)

	log.Print("Out: ", player.QuestState)

	q.playerRepo.Add(player)

	return triggerResponse, nil
}
