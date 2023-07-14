package web

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatteryLevelFromVoltage(t *testing.T) {
	val := batteryLevelFromVoltage(3.5)
	require.Equal(t, uint8(100), val)
}
