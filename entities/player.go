package entities

// Player defines a player for our application
type Player struct {
	ID           uint   `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	QuestID      uint   `json:"questid"`
	QuestState   string `josn:"queststate"`
	Achievements string `json:"achievements"`
	Isnabled     bool   `json:"isenabled"`
}
