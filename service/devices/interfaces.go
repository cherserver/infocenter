package devices

import "time"

type Device interface {
	SID() string
	LastUpdateAt() time.Time
}

type BatteryPowered interface {
	BatteryVoltage() float32
}

type Thermometer interface {
	Temperature() float32
}

type Hygrometer interface {
	Humidity() float32
}

type Barometer interface {
	Pressure() float32 // in Pascals
}

type Sensors interface {
	Sensors() []interface{}
}
