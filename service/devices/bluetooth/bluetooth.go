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
	miDeviceName      = "MJ_HT_V1"

	dev1 = "4c:65:a8:d1:e7:80"
	dev2 = "4c:65:a8:d6:89:74"
)

func NewServer() *Server {
	return &Server{
		ctx: context.Background(),
	}
}

type Server struct {
	ctx context.Context
}

func (s *Server) Init() error {
	device, err := linux.NewDeviceWithName(selfBleDeviceName)
	if err != nil {
		return fmt.Errorf("failed to create default device: %w", err)
	}
	ble.SetDefaultDevice(device)

	err = ble.Scan(s.ctx, false, s.advHandler, nil)
	if err != nil {
		return fmt.Errorf("failed to start bluetooth scanning: %w", err)
	}

	dev1Addr := ble.NewAddr(dev1)
	client, err := ble.Dial(s.ctx, dev1Addr)
	if err != nil {
		return fmt.Errorf("failed to dial dev1: %w", err)
	}

	profile, err := client.DiscoverProfile(true)
	for _, service := range profile.Services {
		log.Printf("Found service: %v", service.UUID.String())
	}

	log.Printf("Blutooth server initialized")

	return nil
}

func (s *Server) Stop() {
}

func (s *Server) advHandler(a ble.Advertisement) {
	if a.LocalName() != miDeviceName {
		return
	}

	addr := a.Addr().String()
	if addr == dev1 || addr == dev2 {
		return
	}

	log.Printf("Found unknown Mi device '%s'", addr)
}
