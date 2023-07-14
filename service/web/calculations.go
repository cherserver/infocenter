package web

const (
	maxVoltage = 3.10
	minVoltage = 2.82

	voltageScale = maxVoltage - minVoltage
)

func batteryLevelFromVoltage(voltage float32) uint8 {
	if voltage >= maxVoltage {
		return 100
	}

	if voltage <= minVoltage {
		return 0
	}

	return uint8(((voltage - minVoltage) / voltageScale) * 100)
}
