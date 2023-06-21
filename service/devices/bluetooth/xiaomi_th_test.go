package bluetooth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleDataNotify_NullTerminated(t *testing.T) {
	var device xiaomiTH
	device.handleDataNotify([]byte("T=27.3 H=30.4\000"))
	require.Equal(t, float32(27.3), device.Temperature())
	require.Equal(t, float32(30.4), device.Humidity())
}

func TestHandleDataNotify_NotNullTerminated(t *testing.T) {
	var device xiaomiTH
	device.handleDataNotify([]byte("T=27.3 H=30.4"))
	require.Equal(t, float32(27.3), device.Temperature())
	require.Equal(t, float32(30.4), device.Humidity())
}
