package repositories

import (
	"log"

	"github.com/go-redis/redis"
)

// StorageManager defines interface
type StorageManager interface {
	Load(key string) (string, error)
	Store(key string, data []byte) error
}

type storageManager struct {
	redisClient *redis.Client
}

// NewStorageManager returns a new UserRepo instance
func NewStorageManager(rc *redis.Client) StorageManager {
	return &storageManager{
		redisClient: rc,
	}
}

// Store saves data to redis
func (s *storageManager) Store(key string, data []byte) error {
	log.Print("Saving JSON")

	log.Print("JSON data: ", string(data))

	return s.redisClient.Set(key, data, 0).Err()
}

// Load retrieves data from redis
func (s *storageManager) Load(key string) (string, error) {
	log.Print("Loading JSON")

	data, err := s.redisClient.Get(key).Result()

	if err != nil {
		return "", err
	}
	log.Print("JSON data: ", data)

	return data, nil
}
