package transport

import (
	"context"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"log"
)

func (t *Transport) processIncomingCmd(msg *message) {
	outChan := t.getAwaitingByName(msg.Cmd)
	if outChan == nil {
		log.Printf("Unexpected response: %+v", msg)
		return
	}

	outChan <- response{
		msg:  msg,
		data: []byte(msg.Data),
		err:  nil,
	}
}

func (t *Transport) queueAwaitingByName(cmd string) (chan response, error) {
	t.awaitingByNameMutex.Lock()
	defer t.awaitingByNameMutex.Unlock()

	if _, fnd := t.awaitingByName[cmd]; fnd {
		return nil, fmt.Errorf("command '%v' is already awaiting response", cmd)
	}

	outChan := make(chan response, 0)
	t.awaitingByName[cmd] = outChan
	return outChan, nil
}

func (t *Transport) getAwaitingByName(cmd string) chan response {
	t.awaitingByNameMutex.Lock()
	defer t.awaitingByNameMutex.Unlock()

	return t.awaitingByName[cmd]
}

func (t *Transport) dropAwaitingByName(cmd string) chan response {
	t.awaitingByNameMutex.Lock()
	defer t.awaitingByNameMutex.Unlock()

	outChan, _ := t.awaitingByName[cmd]
	delete(t.awaitingByName, cmd)

	return outChan
}

func (t *Transport) request(sid string, cmd string, data map[string]interface{}) <-chan response {
	outChan, err := t.queueAwaitingByName(cmd)
	if err != nil {
		outErrChan := make(chan response, 0)

		go func() {
			outErrChan <- response{
				data: nil,
				err:  err,
			}
			close(outErrChan)
		}()

		return outErrChan
	}

	responseChan := make(chan response, 0)

	mode := cipher.NewCBCEncrypter(t.token, initVector)
	cipherText := make([]byte, len(t.gatewayToken))
	mode.CryptBlocks(cipherText, []byte(t.gatewayToken))
	cmdObj := &message{
		Sid: sid,
		Cmd: cmd,
	}

	if data == nil {
		data = make(map[string]interface{}, 1)
	}
	data["key"] = fmt.Sprintf("%X", cipherText)
	bytes, _ := json.Marshal(data)
	cmdObj.Data = string(bytes)

	cmdJson, err := json.Marshal(cmdObj)
	if err != nil {
		go func() {
			responseChan <- response{
				data: nil,
				err:  fmt.Errorf("failed to marshal request: %w", err),
			}
			close(responseChan)
		}()
		return responseChan
	}

	err = t.sendMessage(cmdJson)
	if err != nil {
		log.Printf("Failed to send CMD: %s", err.Error())
		go func() {
			responseChan <- response{
				data: nil,
				err:  fmt.Errorf("failed to send command '%v': %w", cmd, err),
			}
			close(responseChan)
		}()
		return responseChan
	}

	go t.awaitResponseByName(cmd, outChan, responseChan)

	return responseChan
}

func (t *Transport) awaitResponseByName(cmd string, outChan chan response, responseChan chan response) {
	requestCtx, cancelFunc := context.WithTimeout(t.ctx, requestTimeout)
	defer cancelFunc()
	select {
	case <-requestCtx.Done():
		t.dropAwaitingByName(cmd)
		responseChan <- response{
			data: nil,
			err:  fmt.Errorf("request timed out"),
		}
	case responseData := <-outChan:
		t.dropAwaitingByName(cmd)
		responseChan <- responseData
	}

	close(responseChan)
}

func (t *Transport) sendMessage(msg []byte) error {
	log.Printf("Sending msg %s", string(msg))
	_, err := t.conn.Write(msg)
	if err != nil {
		log.Printf("Error writing to UDP: %s", err.Error())
		return err
	}

	return nil
}
