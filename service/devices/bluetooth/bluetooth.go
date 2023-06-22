package bluetooth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

const (
	selfBleDeviceName = "infocenter"
)

func NewServer() *Server {
	return &Server{
		ctx:          context.Background(),
		knownDevices: make(map[string]Device, 0),
	}
}

type Device interface {
	Address() ble.Addr
	Connect()
}

type Server struct {
	ctx          context.Context
	knownDevices map[string]Device
}

func (s *Server) Init() error {
	device, err := linux.NewDeviceWithName(selfBleDeviceName)
	if err != nil {
		return fmt.Errorf("failed to create default device: %w", err)
	}
	ble.SetDefaultDevice(device)

	_ = s.addKnownDevice(newXiaomiTH(dev1))
	_ = s.addKnownDevice(newXiaomiTH(dev2))

	scanCtx, cancelScan := context.WithCancel(s.ctx)
	go func() {
		err = ble.Scan(scanCtx, false, s.advHandler, nil)
		if err != nil {
			log.Printf("Failed to start bluetooth scanning: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)
	cancelScan()

	for _, dev := range s.knownDevices {
		go dev.Connect()
	}

	log.Printf("Blutooth server initialized")

	return nil
}

func (s *Server) Stop() {
	//_ = ble.Stop()
}

func (s *Server) addKnownDevice(device Device) error {
	_, fnd := s.knownDevices[device.Address().String()]
	if fnd {
		return fmt.Errorf("device '%v' is already added in the known devices list", device.Address())
	}

	s.knownDevices[device.Address().String()] = device
	return nil
}

func (s *Server) advHandler(a ble.Advertisement) {
	if a.LocalName() != miDeviceName {
		return
	}

	if device, fnd := s.knownDevices[a.Addr().String()]; fnd {
		log.Printf("Device found: %v", device.Address())
		return
	}

	log.Printf("Found unknown Mi device '%v'", a.Addr())
}
