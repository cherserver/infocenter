package transport

type messageID int64

type message struct {
	Sid   string `json:"sid,omitempty"`
	Model string `json:"model,omitempty"`
	Data  string `json:"data,omitempty"`
	Token string `json:"token,omitempty"`
	Cmd   string `json:"cmd"`
}

type deviceCommand struct {
	ID     messageID     `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params,omitempty"`
}

type response struct {
	msg  *message
	data []byte
	err  error
}
