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
		s.connectDevice(a.Addr())
		return
	}

	log.Printf("Found unknown Mi device '%s'", addr)
}

func (s *Server) connectDevice(addr ble.Addr) {
	client, err := ble.Dial(s.ctx, addr)
	if err != nil {
		log.Fatalf("failed to dial dev1: %v", err)
		return
	}

	profile, err := client.DiscoverProfile(true)
	for _, service := range profile.Services {
		log.Printf("Dev '%v', found service: %v", addr, service.UUID.String())
	}
}
