package bluetooth

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

const (
	selfBleDeviceName = "infocenter"
	miDeviceName      = "MJ_HT_V1"

	dev1 = "4c:65:a8:d1:e7:80"
	dev2 = "4c:65:a8:d6:89:74"

	batteryService = "180f"
	batteryChar    = "2a19"
	dataService    = "226c000064764566756266734470666d"
	dataChar       = "226caa5564764566756266734470666d"
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
		log.Fatalf("failed to dial %v: %v", addr, err)
		return
	}

	profile, err := client.DiscoverProfile(true)
	for _, service := range profile.Services {
		switch service.UUID.String() {
		case batteryService:
			for _, char := range service.Characteristics {
				if char.UUID.String() != batteryChar {
					continue
				}

				if (char.Property & ble.CharRead) != 0 {
					val, err := client.ReadCharacteristic(char)
					if err != nil {
						log.Printf("Failed to read characteristic: %s", err)
						continue
					}

					log.Printf("Device '%v' battery level %d", addr, val)
				}
			}
		case dataService:
			for _, char := range service.Characteristics {
				if char.UUID.String() != dataChar {
					continue
				}

				if (char.Property & ble.CharNotify) != 0 {
					log.Printf("Device '%v' data notify ok", addr)
				}

				err = client.Subscribe(char, false, s.handleDataNotify)
				if err != nil {
					log.Printf("Failed to subscribe to characteristic: %s", err)
					continue
				}
			}
		}
	}
}

func (s *Server) handleDataNotify(req []byte) {
	if len(req) == 0 {
		return
	}

	data := string(req)
	values := strings.Split(data, " ")
	if len(values) != 2 {
		return
	}

	temp := 0.0
	hum := 0.0

	for _, val := range values {
		if len(val) == 0 {
			continue
		}

		switch val[0] {
		case 'T':
			tempVal, err := strconv.ParseFloat(val[1:], 32)
			if err != nil {
				temp = tempVal
			}
		case 'H':
			humVal, err := strconv.ParseFloat(val[1:], 32)
			if err != nil {
				hum = humVal
			}
		}
	}

	log.Printf("Notify data: %s", string(req))
	log.Printf("Temperature: %v, humidity: %v", temp, hum)
}
