package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/cherserver/infocenter/service/devices"
)

type Server struct {
	currentSessionId uuid.UUID
	sensors          devices.Sensors
	listener         net.Listener
}

func NewServer(sensors devices.Sensors) *Server {
	return &Server{
		currentSessionId: uuid.New(),
		sensors:          sensors,
	}
}

func (s *Server) Init() error {
	fileServer := http.FileServer(http.Dir("./http"))
	http.Handle("/", fileServer)

	http.HandleFunc("/reset", s.resetHandler)

	http.HandleFunc("/sensors", s.sensorsHandler)

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
