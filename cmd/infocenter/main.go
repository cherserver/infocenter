package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cherserver/infocenter/service/devices/bluetooth"
)

func main() {
	var err error

	bluetoothServer := bluetooth.NewServer()
	err = bluetoothServer.Init()
	if err != nil {
		log.Fatalf("Failed to initialize bluetooth server: %v", err)
	}

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	stopSignal := <-stopSignalCh
	log.Printf("Signal '%+v' caught, exit", stopSignal)

	bluetoothServer.Stop()
}
