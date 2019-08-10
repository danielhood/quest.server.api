package quests

// FindAllTreasureQuest provides logic for TreasureQuest
type FindAllTreasureQuest interface {
	Trigger(int, string) (string, error)
}

type findAllTreasureQuest struct {
}

// NewFindAllTreasureQuest creates a new FindAllTreasureQuest
func NewFindAllTreasureQuest() FindAllTreasureQuest {
	return &findAllTreasureQuest{}
}

func (q *findAllTreasureQuest) Trigger(playerCode int, deviceType string) (string, error) {
	// Valid action codes are: NONE, UNKNOWN_PLAYER, NO_QUEST, COMPLETED, ACTIVATE

	// TODO: Process game state for player on current quest and return device action code

	return "ACTIVATE", nil
}
