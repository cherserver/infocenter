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
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) Init() error {
	device, err := linux.NewDeviceWithName(selfBleDeviceName)
	if err != nil {
		return fmt.Errorf("failed to create default device: %w", err)
	}
	ble.SetDefaultDevice(device)

	err = ble.Scan(context.Background(), false, s.advHandler, nil)
	if err != nil {
		return fmt.Errorf("failed to start bluetooth scanning: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
}

func (s *Server) advHandler(a ble.Advertisement) {
	if a.LocalName() != miDeviceName {
		return
	}

	log.Printf("Mi temperature device found: %v, %v", a.Addr(), a.Connectable())
}
