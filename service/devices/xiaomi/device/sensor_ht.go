package device

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cherserver/infocenter/service/devices"
	"github.com/cherserver/infocenter/service/devices/xiaomi/transport"
)

var (
	_ devices.Device         = &SensorHT{}
	_ devices.BatteryPowered = &SensorHT{}
	_ devices.Thermometer    = &SensorHT{}
	_ devices.Hygrometer     = &SensorHT{}
)

type sensorHTData struct {
	Voltage     *int    `json:"voltage"`
	Temperature *string `json:"temperature"`
	Humidity    *string `json:"humidity"`
}

type SensorHT struct {
	sid          string
	lastUpdateAt atomic.Pointer[time.Time]
	voltage      atomic.Pointer[float32]
	temperature  atomic.Pointer[float32]
	humidity     atomic.Pointer[float32]
}

func (s *SensorHT) SID() string {
	return s.sid
}

func (s *SensorHT) LastUpdateAt() time.Time {
	ptr := s.lastUpdateAt.Load()
	if ptr == nil {
		return time.Time{}
	}

	return *ptr
}

func (s *SensorHT) OnHeartBeat(data string) {
	_ = s.parseData(data)
}

func (s *SensorHT) OnReport(data string) {
	_ = s.parseData(data)
}

func NewSensorHT(gateway *transport.Transport, sid string, initData string) (*SensorHT, error) {
	sensor := &SensorHT{
		sid: sid,
	}

	var zero float32 = 0
	sensor.voltage.Store(&zero)
	sensor.temperature.Store(&zero)
	sensor.humidity.Store(&zero)

	err := sensor.parseData(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to create sensor_ht with sid '%v', can't parse init data: %w", sid, err)
	}

	gateway.RegisterHeartBeatConsumer(sensor.SID(), sensor.OnHeartBeat)
	gateway.RegisterReportConsumer(sensor.SID(), sensor.OnReport)

	return sensor, nil
}

func (s *SensorHT) BatteryVoltage() float32 {
	return *s.voltage.Load()
}

func (s *SensorHT) Temperature() float32 {
	return *s.temperature.Load()
}

func (s *SensorHT) Humidity() float32 {
	return *s.humidity.Load()
}

func (s *SensorHT) parseData(data string) error {
	var parsedData sensorHTData
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	if parsedData.Voltage != nil {
		voltage := float32(*parsedData.Voltage) / 1000
		s.voltage.Store(&voltage)
	}

	if parsedData.Temperature != nil {
		temperature, err := strconv.ParseInt(*parsedData.Temperature, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse temperature from '%v': %w", *parsedData.Temperature, err)
		}

		floatTemp := float32(temperature) / 100
		s.temperature.Store(&floatTemp)
	}

	if parsedData.Humidity != nil {
		humidity, err := strconv.ParseInt(*parsedData.Humidity, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse humidity from '%v': %w", *parsedData.Humidity, err)
		}

		floatHum := float32(humidity) / 100
		s.humidity.Store(&floatHum)
	}

	log.Printf("New '%s' device data: temp '%v', hum '%v', voltage '%v'",
		s.sid, s.Temperature(), s.Humidity(), s.BatteryVoltage())

	newLastUpdateAt := time.Now()
	s.lastUpdateAt.Store(&newLastUpdateAt)

	return nil
}
