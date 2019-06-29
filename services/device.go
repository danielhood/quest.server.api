package services

import (
	"github.com/danielhood/quest.server.api/entities"
	"github.com/danielhood/quest.server.api/repositories"
)

// DeviceService provides a CRUD interface for Devices
type DeviceService interface {
	ReadAll() ([]entities.Device, error)
	Read(string, string) (*entities.Device, error)
	Update(*entities.Device) error
	Delete(*entities.Device) error
}

// NewDeviceService creates a new DeviceService
func NewDeviceService(dr repositories.DeviceRepo) DeviceService {
	return &deviceService{
		deviceRepo: dr,
	}
}

type deviceService struct {
	deviceRepo repositories.DeviceRepo
}

func (s *deviceService) ReadAll() ([]entities.Device, error) {
	return s.deviceRepo.GetAll()
}

func (s *deviceService) Read(hostname string, deviceKey string) (*entities.Device, error) {
	return s.deviceRepo.GetByHostnameAndKey(hostname, deviceKey)
}

func (s *deviceService) Update(u *entities.Device) error {
	return s.deviceRepo.Add(u)
}

func (s *deviceService) Delete(u *entities.Device) error {
	return s.deviceRepo.Delete(u)
}
