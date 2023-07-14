package transport

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	port             = 9898
	multicastAddress = "224.0.0.50"

	cmdAckSuffix = "_ack"

	requestTimeout = 2 * time.Second

	requestGetChildDevices = "get_id_list"
	requestReadDevice      = "read"

	eventHeartbeat = "heartbeat"
	eventReport    = "report"
)

var (
	initVector = []byte{0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58, 0x56, 0x2e}
)

func New(addr string, token string) (*Transport, error) {
	tokenDecoded, err := aes.NewCipher([]byte(token))
	if err != nil {
		return nil, fmt.Errorf("failed to decode gateway token: %w", err)
	}

	addressDecoded := net.ParseIP(addr)
	if addressDecoded == nil {
		return nil, fmt.Errorf("failed to decode gateway IP-address: %v", addr)
	}

	return &Transport{
		address:        addressDecoded,
		token:          tokenDecoded,
		conn:           nil,
		awaiting:       make(map[messageID]chan response, 0),
		awaitingByName: make(map[string]chan response, 0),
		stopped:        make(chan struct{}),
		ctx:            context.Background(),

		heartBeatConsumers: make(map[string]EventConsumeFunc, 0),
		reportConsumers:    make(map[string]EventConsumeFunc, 0),
	}, nil
}

type Transport struct {
	address       net.IP
	token         cipher.Block
	conn          *net.UDPConn
	multicastConn *net.UDPConn

	gatewaySID   string
	gatewayToken string

	awaiting      map[messageID]chan response
	awaitingMutex sync.Mutex

	awaitingByName      map[string]chan response
	awaitingByNameMutex sync.Mutex

	heartBeatConsumers      map[string]EventConsumeFunc
	heartBeatConsumersMutex sync.Mutex

	reportConsumers      map[string]EventConsumeFunc
	reportConsumersMutex sync.Mutex

	stopped chan struct{}

	ctx context.Context
}

func (t *Transport) Start() error {
	udpAddr := &net.UDPAddr{
		IP:   t.address,
		Port: port,
	}

	var err error
	t.conn, err = net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("failed to dial gateway UDP: %w", err)
	}

	multicastAddr := &net.UDPAddr{
		IP:   net.ParseIP(multicastAddress),
		Port: port,
	}
	t.multicastConn, err = net.ListenMulticastUDP("udp4", nil, multicastAddr)
	if err != nil {
		return fmt.Errorf("failed to dial multicast UDP: %w", err)
	}

	go t.reader()
	go t.multicastReader()

	resp := <-t.request("", requestGetChildDevices, nil)
	if resp.err != nil {
		return fmt.Errorf("failed to call gateway devices list: %w", err)
	}

	t.gatewaySID = resp.msg.Sid
	t.gatewayToken = resp.msg.Token

	t.RegisterHeartBeatConsumer(t.gatewaySID, t.onHeartBeat)

	log.Printf("Transport successfully started (%v)", t.address.String())

	return nil
}

func (t *Transport) Stop() error {
	close(t.stopped)

	if t.conn != nil {
		_ = t.conn.Close()
	}

	return nil
}

func (t *Transport) reader() {
	buf := make([]byte, 2048)
	for {
		select {
		case <-t.stopped:
			log.Printf("Stopped listening for UDP data")
			return
		default:
			size, _, err := t.conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			if size > 0 {
				msg := make([]byte, size)
				copy(msg, buf[0:size])
				go t.processIncomingMessage(msg)
			}
		}
	}
}

func (t *Transport) multicastReader() {
	buf := make([]byte, 2048)
	for {
		select {
		case <-t.stopped:
			log.Printf("Stopped listening for multicast UDP data")
			return
		default:
			size, _, err := t.multicastConn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading from multicast UDP: %v", err)
				continue
			}

			if size > 0 {
				msg := make([]byte, size)
				copy(msg, buf[0:size])
				go t.processIncomingMessage(msg)
			}
		}
	}
}

func (t *Transport) processIncomingMessage(msg []byte) {
	log.Printf("Received UDP message: '%s'", string(msg))
	cmd := &message{}
	err := json.Unmarshal(msg, cmd)
	if err != nil {
		log.Printf("Failed to unmarshal command: %v", err)
		return
	}

	messageName, isCommand := strings.CutSuffix(cmd.Cmd, cmdAckSuffix)
	cmd.Cmd = messageName
	if isCommand {
		t.processIncomingCmd(cmd)
	} else {
		t.processIncomingEvent(cmd)
	}
}

func (t *Transport) onHeartBeat(string) {
	log.Printf("Got gateway heartbeat")
}
