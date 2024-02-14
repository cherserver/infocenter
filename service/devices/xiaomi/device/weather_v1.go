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
	_ devices.Device         = &WeatherV1{}
	_ devices.BatteryPowered = &WeatherV1{}
	_ devices.Thermometer    = &WeatherV1{}
	_ devices.Hygrometer     = &WeatherV1{}
	_ devices.Barometer      = &WeatherV1{}
)

type weatherV1Data struct {
	Voltage     *int    `json:"voltage"`
	Temperature *string `json:"temperature"`
	Humidity    *string `json:"humidity"`
	Pressure    *string `json:"pressure"`
}

type WeatherV1 struct {
	sid          string
	lastUpdateAt atomic.Pointer[time.Time]
	voltage      atomic.Pointer[float32]
	temperature  atomic.Pointer[float32]
	humidity     atomic.Pointer[float32]
	pressure     atomic.Pointer[float32]
}

func (s *WeatherV1) SID() string {
	return s.sid
}

func (s *WeatherV1) LastUpdateAt() time.Time {
	ptr := s.lastUpdateAt.Load()
	if ptr == nil {
		return time.Time{}
	}

	return *ptr
}

func (s *WeatherV1) OnHeartBeat(data string) {
	_ = s.parseData(data)
}

func (s *WeatherV1) OnReport(data string) {
	_ = s.parseData(data)
}

func NewWeatherV1(gateway *transport.Transport, sid string, initData string) (*WeatherV1, error) {
	sensor := &WeatherV1{
		sid: sid,
	}

	var zero float32 = 0
	sensor.voltage.Store(&zero)
	sensor.temperature.Store(&zero)
	sensor.pressure.Store(&zero)
	sensor.humidity.Store(&zero)

	err := sensor.parseData(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to create sensor_ht with sid '%v', can't parse init data: %w", sid, err)
	}

	gateway.RegisterHeartBeatConsumer(sensor.SID(), sensor.OnHeartBeat)
	gateway.RegisterReportConsumer(sensor.SID(), sensor.OnReport)

	return sensor, nil
}

func (s *WeatherV1) BatteryVoltage() float32 {
	return *s.voltage.Load()
}

func (s *WeatherV1) Temperature() float32 {
	return *s.temperature.Load()
}

func (s *WeatherV1) Humidity() float32 {
	return *s.humidity.Load()
}

func (s *WeatherV1) Pressure() float32 {
	return *s.pressure.Load()
}

func (s *WeatherV1) parseData(data string) error {
	var parsedData weatherV1Data
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

	if parsedData.Pressure != nil {
		pressure, err := strconv.ParseInt(*parsedData.Pressure, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse pressure from '%v': %w", *parsedData.Pressure, err)
		}

		floatPressure := float32(pressure)
		s.pressure.Store(&floatPressure)
	}

	log.Printf("New '%s' device data: temp '%v', hum '%v', pressure '%v', voltage '%v'",
		s.sid, s.Temperature(), s.Humidity(), s.Pressure(), s.BatteryVoltage())

	newLastUpdateAt := time.Now()
	s.lastUpdateAt.Store(&newLastUpdateAt)

	return nil
}
