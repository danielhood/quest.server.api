package entities

// Quest defines a quest entity for our application
type Quest struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"password"`
	Key       string `josn:"key"`
	IsEnabled bool   `json:"isenabled"`
}
