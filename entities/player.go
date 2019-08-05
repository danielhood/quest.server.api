package entities

// Player defines a player for our application
type Player struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	QuestKey     string `json:"questkey"`
	QuestState   string `json:"queststate"`
	Achievements string `json:"achievements"`
	IsEnabled    bool   `json:"isenabled"`
}
