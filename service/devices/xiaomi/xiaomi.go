package xiaomi

import (
	"fmt"
	"log"

	"github.com/cherserver/infocenter/service/devices"
	"github.com/cherserver/infocenter/service/devices/xiaomi/device"
	"github.com/cherserver/infocenter/service/devices/xiaomi/transport"
)

var _ devices.Sensors = &Gateway{}

func NewGateway(addr string, token string) (*Gateway, error) {
	trans, err := transport.New(addr, token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode gateway token: %w", err)
	}

	return &Gateway{
		transport: trans,
	}, nil
}

type Gateway struct {
	transport    *transport.Transport
	childDevices []interface{}
}

func (g *Gateway) Sensors() []interface{} {
	return g.childDevices
}

func (g *Gateway) Init() error {
	err := g.transport.Start()
	if err != nil {
		return fmt.Errorf("failed to start transpot: %w", err)
	}

	devices, err := g.transport.RequestGetChildDevicesIDs()
	if err != nil {
		return fmt.Errorf("failed to get devices: %w", err)
	}

	log.Printf("Devices: %v", devices)
	for _, deviceSID := range devices {
		deviceInfo, err := g.transport.RequestReadDevice(deviceSID)
		if err != nil {
			return fmt.Errorf("failed to read device '%v' info: %w", deviceSID, err)
		}

		switch deviceInfo.Model {
		case "sensor_ht":
			sensor, err := device.NewSensorHT(g.transport, deviceSID, deviceInfo.Data)
			if err != nil {
				return fmt.Errorf("failed to create sensor_ht: %w", err)
			}

			g.childDevices = append(g.childDevices, sensor)
		case "weather.v1":
			sensor, err := device.NewWeatherV1(g.transport, deviceSID, deviceInfo.Data)
			if err != nil {
				return fmt.Errorf("failed to create weather.v1: %w", err)
			}

			g.childDevices = append(g.childDevices, sensor)
		}
	}

	log.Printf("Gateway successfully started")

	return nil
}

func (g *Gateway) Stop() {
	_ = g.transport.Stop()
}
