package entities

// Device defines a user for our application
type Device struct {
	Hostname     string `json:"hostname"`
	DeviceKey    string `json:"devicekey"`
	IsRegistered bool   `json:"isregistered"`
	IsEnabled    bool   `json:"isenabled"`
}
