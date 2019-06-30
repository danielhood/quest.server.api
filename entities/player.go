package entities

// Player defines a player for our application
type Player struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	QuestKey     uint   `json:"questkey"`
	QuestState   string `josn:"queststate"`
	Achievements string `json:"achievements"`
	Isnabled     bool   `json:"isenabled"`
}
