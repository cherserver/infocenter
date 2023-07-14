package transport

import (
	"encoding/json"
	"fmt"
)

func (t *Transport) RequestGetChildDevicesIDs() ([]string, error) {
	resp := <-t.request(t.gatewaySID, requestGetChildDevices, nil)
	if resp.err != nil {
		return nil, fmt.Errorf("failed to get devices list: %w", resp.err)
	}

	devList := make([]string, 0)
	err := json.Unmarshal(resp.data, &devList)

	if err != nil {
		return nil, fmt.Errorf("failed to parse get devices list response: %w", err)
	}

	return devList, nil
}
