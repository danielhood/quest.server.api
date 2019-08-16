package quests

// QuestResponseNone - no resonse (do we need this?)
const QuestResponseNone string = "NONE"

// QuestResponseCompleted - Quest has been completed
const QuestResponseCompleted string = "COMPLETED"

// QuestResponseItemNotPartOfQuest - this item is not part of the quest
const QuestResponseItemNotPartOfQuest string = "ITEM_NOT_PART_OF_QUEST"

// QuestResponesItemAlreadyCollected - this item has already been collected (and quest not completed)
const QuestResponesItemAlreadyCollected string = "ITEM_ALREADY_COLLECTED"

// QuestResponseActivate - item has been collected, activate
const QuestResponseActivate string = "ACTIVATE"

// QuestResponseUnknownPlayer - triggering player is not known to the system
const QuestResponseUnknownPlayer string = "UNKNOWN_PLAYER"

// QuestResponseNoQuest - player does not have an active quest
const QuestResponseNoQuest string = "NO_QUEST"

// QuestResponseUnknownQuest - current quest isn't known to the system
const QuestResponseUnknownQuest string = "UNKNOWN_QUEST"
