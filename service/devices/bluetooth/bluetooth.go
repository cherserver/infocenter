package bluetooth

import (
	"context"
	"fmt"
	"log"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

const (
	selfBleDeviceName = "infocenter"
)

func NewServer() *Server {
	return &Server{
		ctx:          context.Background(),
		knownDevices: make(map[ble.Addr]Device, 0),
	}
}

type Device interface {
	Address() ble.Addr
	Connect()
}

type Server struct {
	ctx          context.Context
	knownDevices map[ble.Addr]Device
}

func (s *Server) Init() error {
	device, err := linux.NewDeviceWithName(selfBleDeviceName)
	if err != nil {
		return fmt.Errorf("failed to create default device: %w", err)
	}
	ble.SetDefaultDevice(device)

	_ = s.addKnownDevice(newXiaomiTH(dev1))
	_ = s.addKnownDevice(newXiaomiTH(dev2))

	err = ble.Scan(s.ctx, false, s.advHandler, nil)
	if err != nil {
		return fmt.Errorf("failed to start bluetooth scanning: %w", err)
	}

	log.Printf("Blutooth server initialized")

	return nil
}

func (s *Server) Stop() {
}

func (s *Server) addKnownDevice(device Device) error {
	_, fnd := s.knownDevices[device.Address()]
	if fnd {
		return fmt.Errorf("device '%v' is already added in the known devices list", device.Address())
	}

	s.knownDevices[device.Address()] = device
	return nil
}

func (s *Server) advHandler(a ble.Advertisement) {
	if a.LocalName() != miDeviceName {
		return
	}

	if device, fnd := s.knownDevices[a.Addr()]; fnd {
		device.Connect()
		return
	}

	log.Printf("Found unknown Mi device '%v'", a.Addr())
}
