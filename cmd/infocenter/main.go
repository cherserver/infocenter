package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cherserver/infocenter/service/devices/xiaomi"
	"github.com/cherserver/infocenter/service/weather"
	"github.com/cherserver/infocenter/service/web"
)

const (
	gatewayAddr  = "192.168.31.21"
	gatewayToken = "540b4bf40bb290ef62004d27fc3438e6"
)

func main() {
	gateway, err := xiaomi.NewGateway(gatewayAddr, gatewayToken)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}

	err = gateway.Init()
	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}

	weatherSource := weather.New()
	err = weatherSource.Init()
	if err != nil {
		log.Fatalf("Failed to initialize weather: %v", err)
	}

	webServer := web.NewServer(gateway, weatherSource)
	err = webServer.Init()
	if err != nil {
		log.Fatalf("Failed to initialize web server: %v", err)
	}

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	stopSignal := <-stopSignalCh
	log.Printf("Signal '%+v' caught, exit", stopSignal)

	webServer.Stop()
	weatherSource.Stop()
	gateway.Stop()
}
