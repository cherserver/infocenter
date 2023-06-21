package devices

type Thermometer interface {
	Temperature() float32
}

type Hygrometer interface {
	Humidity() float32
}

type Battery interface {
	RemainingPowerPercent() uint32
}
