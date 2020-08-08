package services

import (
	"bytes"

	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// AudioService provides a CRUD interface for Audio records
type AudioService interface {
	ReadAll() ([]entities.Audio, error)
	Read(string, string) (*entities.Audio, error)
	Update(*entities.Audio, *bytes.Buffer) error
	Delete(*entities.Audio) error
}

// NewAudioService creates a new AudioService
func NewAudioService(ar repositories.AudioRepo) AudioService {
	return &audioService{
		audioRepo: ar,
	}
}

type audioService struct {
	audioRepo repositories.AudioRepo
}

func (s *audioService) ReadAll() ([]entities.Audio, error) {
	return s.audioRepo.GetAll()
}

func (s *audioService) Read(deviceType string, responseType string) (*entities.Audio, error) {
	return s.audioRepo.GetByDeviceAndResponseTypes(deviceType, responseType)
}

func (s *audioService) Update(u *entities.Audio, fileBuffer *bytes.Buffer) error {
	return s.audioRepo.Add(u, fileBuffer)
}

func (s *audioService) Delete(u *entities.Audio) error {
	return s.audioRepo.Delete(u)
}
