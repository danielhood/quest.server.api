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

// QuestKey holds string key for this quest
const questKey string = "FIND_ALL_TREASURE"

// NewFindAllTreasureQuest creates a new FindAllTreasureQuest
func NewFindAllTreasureQuest(pr repositories.PlayerRepo) FindAllTreasureQuest {
	return &findAllTreasureQuest{
		playerRepo: pr,
	}
}

func (q *findAllTreasureQuest) hasCompletedQuest(player *entities.Player) bool {
	for _, key := range player.Achievements {
		if key == questKey {
			return true
		}
	}
	return false
}

func (q *findAllTreasureQuest) Trigger(player *entities.Player, deviceType string) (string, error) {
	if q.hasCompletedQuest(player) {
		return QuestResponseCompleted, nil
	}

	if player.QuestStatus == "" {
		player.QuestStatus = QuestStatusActive
		player.QuestState, _ = json.Marshal(&findAllTreasureQuestState{false, false, false, false})
	}

	var questState findAllTreasureQuestState
	json.Unmarshal(player.QuestState, &questState)

	triggerResponse := QuestResponseActivate

	switch deviceType {
	case DeviceTypeTreasure1:
		questState.hasTreasure1 = true
		break
	case DeviceTypeTreasure2:
		questState.hasTreasure2 = true
		break
	case DeviceTypeTreasure3:
		questState.hasTreasure3 = true
		break
	case DeviceTypeTreasure4:
		questState.hasTreasure4 = true
		break
	default:
		triggerResponse = QuestResponseWrongItem
	}

	if questState.hasTreasure1 && questState.hasTreasure2 && questState.hasTreasure3 && questState.hasTreasure4 {
		player.QuestStatus = QuestStatusCompleted
		log.Print("QuestStatus: ", QuestStatusCompleted)
		player.Achievements = append(player.Achievements, questKey)
	}

	player.QuestState, _ = json.Marshal(questState)

	q.playerRepo.Add(player)

	return triggerResponse, nil
}
