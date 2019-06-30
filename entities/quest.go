package entities

// Quest defines a quest entity for our application
type Quest struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Key       string `josn:"key"`
	IsEnabled bool   `json:"isenabled"`
}
