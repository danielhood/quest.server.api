package services

type Ping interface {
  Get() string
}

func NewPing() Ping {
  return &PingService{}
}

type PingService struct {}

func(* PingService) Get() string {
  return "pong"
}
