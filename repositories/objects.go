package repositories

import (
  "errors"
  "log"

  "github.com/danielhood/loco.server/entities"
)

// TODO: Move this to somewhere better (redis?)
var objects map[uint]*entities.Object

func init() {
  objects = make(map[uint]*entities.Object)
}

type ObjectRepo interface {
  GetAll() ([]entities.Object, error)
  Get(id uint) (*entities.Object, error)
  Add(o *entities.Object) error
}

type objectRepo struct {
}

func NewObjectRepo() ObjectRepo {
  return &objectRepo{}
}

func (r *objectRepo)GetAll() ([]entities.Object, error) {
  allObjects := make([]entities.Object, len(objects))
  
  idx := 0
  for _, value := range objects {
    allObjects[idx] = *value
    idx++
  }

  return allObjects, nil
}

func (r *objectRepo)Get(id uint) (*entities.Object, error) {
  if val, ok := objects[id]; ok {
    return val, nil;
  }

  return nil, errors.New("Object for id not found")
}

func (r *objectRepo)Add(o *entities.Object) error {
  log.Print("Add object: ", o.Name)

  existing, _ := r.Get(o.Id)
  if existing != nil {
    return errors.New("Object already exists")  // TODO: we could do a merge here
  }

  objects[o.Id] = o

  return nil
}
