package entities

const (
  ObjectType1 = iota
  ObjectType2 = iota
  ObjectType3 = iota
)

type Object struct {
  Id uint `json:"id"`
  Type uint `json:"type"`
  Name string `json:"name"`
  X float32 `json:"x"`
  Y float32 `json:"y"`
  Bearing float32  `json:"b"` // radains
  Velocity float32 `json:"v"` // units??
}
