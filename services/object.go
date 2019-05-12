package services

import (
  "encoding/json"

  "github.com/danielhood/quest.server.api/repositories"
)

type Object interface {
  GetAll() ([]byte, error)
}

func NewObject() Object {
  return &ObjectService{
    repositories.NewObjectRepo(),
  }
}

type ObjectService struct {
  objectRepo repositories.ObjectRepo
}

func(o *ObjectService) GetAll() ([]byte, error) {
  objects, err := o.objectRepo.GetAll()
  if (err == nil) {
    return json.Marshal(objects)
  }

  return nil, err
}
