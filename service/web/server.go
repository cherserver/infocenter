package web

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cherserver/infocenter/service/devices"
	"github.com/cherserver/infocenter/service/weather"
)

const (
	weatherImgPrefix = "/img/weather/64x64/"
)

type Server struct {
	currentSessionId uuid.UUID

	sensors       devices.Sensors
	weatherSource weather.Info

	listener net.Listener
}

func NewServer(sensors devices.Sensors, weatherSource weather.Info) *Server {
	return &Server{
		currentSessionId: uuid.New(),
		sensors:          sensors,
		weatherSource:    weatherSource,
	}
}

func (s *Server) Init() error {
	fileServer := http.FileServer(http.Dir("./http"))
	http.Handle("/", fileServer)

	http.HandleFunc("/reset", s.resetHandler)

	http.HandleFunc("/sensors", s.sensorsHandler)
	http.HandleFunc("/weather", s.weatherHandler)

	server := &http.Server{Addr: ":80", Handler: nil}
	var err error
	s.listener, err = net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	go func() {
		srvError := server.Serve(s.listener)
		log.Printf("HTTP server stopped: %v", srvError)
	}()

	log.Printf("Web server started")
	return nil
}

func (s *Server) Stop() {
	_ = s.listener.Close()
	log.Printf("Web server stopped")
}

func (s *Server) resetHandler(http.ResponseWriter, *http.Request) {
	log.Printf("reset caught, exit")
	os.Exit(0)
}

func (s *Server) sensorsHandler(w http.ResponseWriter, r *http.Request) {
	_ = r

	allDevices := s.sensors.Sensors()

	sensors := make([]Sensor, 0, len(allDevices))
	for _, data := range allDevices {
		device, ok := data.(devices.Device)
		if !ok {
			continue
		}

		sensor := Sensor{SID: device.SID()}
		s.fillUpSensor(data, &sensor)

		if !device.LastUpdateAt().IsZero() {
			diff := uint64(time.Now().Sub(device.LastUpdateAt()).Seconds())
			sensor.LastUpdateSec = &diff
		}

		sensors = append(sensors, sensor)
	}

	statusData, err := json.Marshal(sensors)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode status: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("status: %s", statusData)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(statusData)
}

func (s *Server) fillUpSensor(data interface{}, sensor *Sensor) {
	if dev, ok := data.(devices.BatteryPowered); ok {
		val := batteryLevelFromVoltage(dev.BatteryVoltage())
		sensor.BatteryPercent = &val
	}

	if dev, ok := data.(devices.Thermometer); ok {
		val := dev.Temperature()
		sensor.Temperature = &val
	}

	if dev, ok := data.(devices.Hygrometer); ok {
		val := dev.Humidity()
		sensor.Humidity = &val
	}

	if dev, ok := data.(devices.Barometer); ok {
		val := dev.Pressure()
		sensor.Pressure = &val
	}
}

func (s *Server) weatherHandler(w http.ResponseWriter, r *http.Request) {
	_ = r

	weatherInfo := &Weather{
		Current:  s.currentWeather(),
		Forecast: s.weatherForecast(),
	}

	weatherResponse, err := json.Marshal(weatherInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode weather: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("weather: %s", weatherResponse)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(weatherResponse)
}

func (s *Server) currentWeather() CurrentWeather {
	currWeather := s.weatherSource.CurrentWeather()

	phase := "day"
	if !currWeather.IsDay {
		phase = "night"
	}

	return CurrentWeather{
		ConditionText:  currWeather.ConditionText,
		ConditionImage: fmt.Sprintf(weatherImgPrefix+"%s/%s.png", phase, currWeather.ConditionImageCode),
		Temperature:    currWeather.Temperature,
		Wind:           currWeather.Wind,
		Gust:           currWeather.Gust,
		WindDegree:     currWeather.WindDegree,
		WindDir:        currWeather.WindDir,
		Pressure:       currWeather.Pressure,
		Precipitation:  currWeather.Precipitation,
		Humidity:       currWeather.Humidity,
		CloudPercent:   currWeather.CloudPercent,
		FeelsLike:      currWeather.FeelsLike,
		Visibility:     currWeather.Visibility,
		UV:             currWeather.UV,
	}
}

func (s *Server) weatherForecast() []ForecastItem {
	forecastData := s.weatherSource.Forecast()
	forecast := make([]ForecastItem, 0, len(forecastData))

	for _, item := range forecastData {
		forecast = append(forecast, ForecastItem{
			ConditionText:      item.ConditionText,
			ConditionImage:     fmt.Sprintf(weatherImgPrefix+"day/%s.png", item.ConditionImageCode),
			MaxTemp:            item.MaxTemp,
			MinTemp:            item.MinTemp,
			AvgTemp:            item.AvgTemp,
			MaxWind:            item.MaxWind,
			TotalPrecipitation: item.TotalPrecipitation,
			TotalSnow:          item.TotalSnow,
			AvgVis:             item.AvgVis,
			AvgHumidity:        item.AvgHumidity,
			DailyWillItRain:    item.DailyWillItRain,
			DailyChanceOfRain:  item.DailyChanceOfRain,
			DailyWillItSnow:    item.DailyWillItSnow,
			DailyChanceOfSnow:  item.DailyChanceOfSnow,
			UV:                 item.UV,
		})
	}

	return forecast
}
