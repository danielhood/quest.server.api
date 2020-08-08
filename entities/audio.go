package entities

// Audio defines an audio record for our application
type Audio struct {
	DeviceType   string `json:"devicetype"`
	ResponseType string `json:"responsetype"`
	FilePath     string `json:"filepath"`
}
