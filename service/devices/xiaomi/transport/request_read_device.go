package transport

import (
	"fmt"
)

type ResponseReadDevice struct {
	Model string
	Data  string
}

func (t *Transport) RequestReadDevice(sid string) (*ResponseReadDevice, error) {
	resp := <-t.request(sid, requestReadDevice, nil)
	if resp.err != nil {
		return nil, fmt.Errorf("failed to get devices list: %w", resp.err)
	}

	return &ResponseReadDevice{
		Model: resp.msg.Model,
		Data:  resp.msg.Data,
	}, nil
}
