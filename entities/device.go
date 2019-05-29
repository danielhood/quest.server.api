package entities

// Device defines a user for our application
type Device struct {
	Hostname     string `json:"hostname"`
	IsRegistered bool   `json:"isregistered"`
	Key          string `josn:"key"`
	IsEnabled    bool   `json:"isenabled"`
}
