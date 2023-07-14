package web

type Sensor struct {
	SID            string   `json:"sid"`
	BatteryPercent *uint8   `json:"battery_percent,omitempty"`
	Temperature    *float32 `json:"temperature,omitempty"`
	Humidity       *float32 `json:"humidity,omitempty"`
	Pressure       *float32 `json:"pressure,omitempty"`
}
