package repositories

// https://github.com/go-redis/redis

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/danielhood/quest.server.api/entities"
)

var devices []entities.Device

func init() {
	// Default strucutre
	devices = make([]entities.Device, 0)
}

// DeviceRepo defines DeviceRepo interface
type DeviceRepo interface {
	GetAll() ([]entities.Device, error)
	GetByHostnameAndKey(hostname string, key string) (*entities.Device, error)
	Add(o *entities.Device) error
	Delete(o *entities.Device) error
}

type deviceRepo struct {
	storageManager StorageManager
}

// NewDeviceRepo returns a new DeviceRepo instance
func NewDeviceRepo(sm StorageManager) DeviceRepo {
	r := deviceRepo{
		storageManager: sm,
	}

	r.load()

	return &r
}

func (r *deviceRepo) GetAll() ([]entities.Device, error) {
	allDevices := make([]entities.Device, len(devices))

	idx := 0
	for _, value := range devices {
		allDevices[idx] = value
		idx++
	}

	return allDevices, nil
}

func (r *deviceRepo) GetByHostnameAndKey(hostname string, deviceKey string) (*entities.Device, error) {
	for i, u := range devices {
		if u.Hostname == hostname && u.DeviceKey == deviceKey {
			return &devices[i], nil
		}
	}

	return nil, errors.New("Device for hostname and key not found")
}

func (r *deviceRepo) Add(d *entities.Device) error {
	log.Print("Add Device: ", d.Hostname, " Key: ", d.DeviceKey)

	existing, _ := r.GetByHostnameAndKey(d.Hostname, d.DeviceKey)
	if existing != nil {
		existing.IsEnabled = d.IsEnabled
		existing.IsRegistered = d.IsRegistered
		existing.DeviceType = d.DeviceType
	} else {
		devices = append(devices, *d)
	}

	return r.store()
}

func (r *deviceRepo) Delete(d *entities.Device) error {
	log.Print("Delete Device: ", d.Hostname, " Key: ", d.DeviceKey)

	for i, device := range devices {
		if device.Hostname == d.Hostname && device.DeviceKey == d.DeviceKey {
			devices[i] = devices[len(devices)-1]
			devices = devices[:len(devices)-1]
			return r.store()
		}
	}

	return nil
}

// Store saves data to redis
func (r *deviceRepo) store() error {
	log.Print("Saving devices")

	devicesJSON, err := json.Marshal(devices)
	if err != nil {
		return err
	}

	return r.storageManager.Store("devices", devicesJSON)
}

// Load retrieves data from redis
func (r *deviceRepo) load() error {
	devicesJSON, err := r.storageManager.Load("devices")

	if err != nil {
		return err
	}

	if len(devicesJSON) == 0 {
		return nil
	}

	if err = json.Unmarshal([]byte(devicesJSON), &devices); err != nil {
		return err
	}

	log.Printf("Loaded %v device(s)", len(devices))

	return nil
}
