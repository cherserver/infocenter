package bluetooth

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/go-ble/ble"

	"github.com/cherserver/infocenter/service/devices"
)

const (
	miDeviceName = "MJ_HT_V1"

	dev1 = "4C:65:A8:D1:E7:80"
	dev2 = "4C:65:A8:D6:89:74"

	batteryService = "180f"
	batteryChar    = "2a19"
	dataService    = "226c000064764566756266734470666d"
	dataChar       = "226caa5564764566756266734470666d"
)

type XiaomiTH interface {
	devices.Battery
	devices.Thermometer
	devices.Hygrometer
}

var _ XiaomiTH = &xiaomiTH{}
var _ Device = &xiaomiTH{}

func newXiaomiTH(address string) *xiaomiTH {
	device := &xiaomiTH{
		ctx:     context.Background(),
		address: ble.NewAddr(address),
	}

	device.batteryLevel.Store(0)
	device.dataPtr.Store(&xiaomiTHData{
		temperature: 0,
		humidity:    0,
	})

	return device
}

type xiaomiTH struct {
	ctx context.Context

	address      ble.Addr
	dataPtr      atomic.Pointer[xiaomiTHData]
	batteryLevel atomic.Uint32
}

type xiaomiTHData struct {
	temperature float32
	humidity    float32
}

func (x *xiaomiTH) Address() ble.Addr {
	return x.address
}

func (x *xiaomiTH) RemainingPowerPercent() uint32 {
	return x.batteryLevel.Load()
}

func (x *xiaomiTH) Temperature() float32 {
	return x.dataPtr.Load().temperature
}

func (x *xiaomiTH) Humidity() float32 {
	return x.dataPtr.Load().humidity
}

func (x *xiaomiTH) Connect() {
	log.Printf("Connect device '%v'", x.address)
	client, err := ble.Dial(x.ctx, x.address)
	if err != nil {
		log.Fatalf("failed to dial %v: %v", x.address, err)
		return
	}

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("failed to discover device profile %v: %v", x.address, err)
		return
	}

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

					log.Printf("Device '%v' battery level %d", x.address, val)
				}
			}
		case dataService:
			for _, char := range service.Characteristics {
				if char.UUID.String() != dataChar {
					continue
				}

				if (char.Property & ble.CharNotify) != 0 {
					log.Printf("Device '%v' data notify ok", x.address)
				}

				err = client.Subscribe(char, false, x.handleDataNotify)
				if err != nil {
					log.Printf("Failed to subscribe to characteristic: %s", err)
					continue
				}
			}
		}
	}

	go func() {
		<-client.Disconnected()
		log.Printf("Clent '%v' disconnected", client.Addr())
		x.Connect()
	}()
}

func (x *xiaomiTH) Disconnect() {

}

func (x *xiaomiTH) handleDataNotify(req []byte) {
	if len(req) == 0 {
		return
	}

	data := ""
	if req[len(req)-1] == 0 {
		data = strings.TrimSpace(string(req[:len(req)-1]))
	} else {
		data = strings.TrimSpace(string(req))
	}

	values := strings.Split(data, " ")
	if len(values) != 2 {
		return
	}

	newData := &xiaomiTHData{
		temperature: 0,
		humidity:    0,
	}

	for _, val := range values {
		if len(val) <= 2 {
			continue
		}

		switch val[0] {
		case 'T':
			tempVal, err := strconv.ParseFloat(val[2:], 32)
			if err == nil {
				newData.temperature = float32(tempVal)
			}
		case 'H':
			humVal, err := strconv.ParseFloat(val[2:], 32)
			if err == nil {
				newData.humidity = float32(humVal)
			}
		}
	}

	log.Printf("Notify data: %s", string(req))
	log.Printf("New data: %v", newData)

	x.dataPtr.Store(newData)
}
