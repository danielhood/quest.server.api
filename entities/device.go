package entities

// Device defines a user for our application
type Device struct {
	ID         uint   `json:"id"`
	Hostname   string `json:"hostname"`
	Registered bool   `json:"registered"`
	Key        string `josn:"key"`
	IsEnabled  bool   `json:"isenabled"`
}
