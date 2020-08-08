package repositories

// https://github.com/go-redis/redis

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

var audioRecords []entities.Audio

func init() {
	// Default strucutre
	audioRecords = make([]entities.Audio, 0)
}

// AudioRepo defines AudioRepo interface
type AudioRepo interface {
	GetAll() ([]entities.Audio, error)
	GetByDeviceAndResponseTypes(deviceType string, responseType string) (*entities.Audio, error)
	Add(o *entities.Audio, fileBuffer *bytes.Buffer) error
	Delete(o *entities.Audio) error
}

type audioRepo struct {
	storageManager StorageManager
}

// NewAudioRepo returns a new AudioRepo instance
func NewAudioRepo(sm StorageManager) AudioRepo {
	r := audioRepo{
		storageManager: sm,
	}

	r.load()

	return &r
}

func (r *audioRepo) GetAll() ([]entities.Audio, error) {
	allAudioRecords := make([]entities.Audio, len(audioRecords))

	idx := 0
	for _, value := range audioRecords {
		allAudioRecords[idx] = value
		idx++
	}

	return allAudioRecords, nil
}

func (r *audioRepo) GetByDeviceAndResponseTypes(deviceType string, responseType string) (*entities.Audio, error) {
	for i, u := range audioRecords {
		if u.DeviceType == deviceType && u.ResponseType == responseType {
			return &audioRecords[i], nil
		}
	}

	return nil, errors.New("Audio for device and response type not found")
}

func (r *audioRepo) Add(d *entities.Audio, fileBuffer *bytes.Buffer) error {
	log.Print("Add Audio: ", d.DeviceType, " - ", d.ResponseType)

	err := r.saveFile(d, fileBuffer)

	if err != nil {
		return err
	}

	existing, _ := r.GetByDeviceAndResponseTypes(d.DeviceType, d.ResponseType)
	if existing != nil {
		existing.FilePath = d.FilePath
	} else {
		audioRecords = append(audioRecords, *d)
	}

	return r.store()
}

func (r *audioRepo) Delete(d *entities.Audio) error {
	log.Print("Delete Audio: ", d.DeviceType, " - ", d.ResponseType)

	for i, audio := range audioRecords {
		if audio.DeviceType == d.DeviceType && audio.ResponseType == d.ResponseType {
			audioRecords[i] = audioRecords[len(audioRecords)-1]
			audioRecords = audioRecords[:len(audioRecords)-1]
			return r.store()
		}
	}

	return nil
}

// Store saves data to redis
func (r *audioRepo) store() error {
	log.Print("Saving audio")

	audioJSON, err := json.Marshal(audioRecords)
	if err != nil {
		return err
	}

	return r.storageManager.Store("audio", audioJSON)
}

// Load retrieves data from redis
func (r *audioRepo) load() error {
	audioJSON, err := r.storageManager.Load("audio")

	if err != nil {
		return err
	}

	if len(audioJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(audioJSON), &audioRecords); err != nil {
		return err
	}

	log.Printf("Loaded %v audio record(s)", len(audioRecords))

	return nil
}

func (r *audioRepo) saveFile(audio *entities.Audio, fileBuffer *bytes.Buffer) error {
	// Compose file based off of fixed folder and type attriutes
	audio.FilePath = "audio/" + audio.DeviceType + "_" + audio.ResponseType

	err := ioutil.WriteFile(audio.FilePath, fileBuffer.Bytes(), 0644)

	if err != nil {
		return err
	}

	return nil
}
