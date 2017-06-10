package entities

const (
  ObjectType1 = "type1"
  ObjectType2 = "type2"
  ObjectType3 = "type3"
)

type Object struct {
  Id uint `json:"id"`
  Type string `json:"type"`
  Name string `json:"name"`
  X float32 `json:"x"`
  Y float32 `json:"y"`
  Bearing float32  `json:"b"` // radains
  Velocity float32 `json:"v"` // units??
}
