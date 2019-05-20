package entities

const (
	ObjectType1 = "type1"
	ObjectType2 = "type2"
	ObjectType3 = "type3"
)

// Object defines general object with location
type Object struct {
	ID       uint    `json:"id"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Bearing  float32 `json:"b"` // radains
	Velocity float32 `json:"v"` // units??
}
