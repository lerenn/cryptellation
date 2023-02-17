package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/pkg/http/health"
	"github.com/digital-feather/cryptellation/services/candlesticks"
)

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init service
	svc, err := candlesticks.New(loadConfigFromEnv())
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initializing application: %w", err))
		return 255
	}
	defer svc.Close()

	// Service marked as ready
	log.Println("Service is ready")
	h.Ready(true)

	// Wait for interrupt
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	return 0
}

func main() {
	os.Exit(run())
}
