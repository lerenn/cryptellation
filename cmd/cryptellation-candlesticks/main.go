package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/digital-feather/cryptellation/internal/candlesticks/app"
	"github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats"
	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/exchanges"
	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/sql"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/health"
)

func initApp() (app.Port, error) {
	// Init database client
	db, err := sql.New(config.LoadSQLConfigFromEnv())
	if err != nil {
		return nil, err
	}

	// Init exchanges connections
	exchanges, err := exchanges.New(config.LoadExchangesConfigFromEnv())
	if err != nil {
		return nil, err
	}

	// Init component
	return app.New(db, exchanges), nil
}

func initController(component app.Port) (func(), error) {
	// Init NATS controller
	natsController, err := nats.NewServer(config.LoadNATSConfigFromEnv(), component)
	if err != nil {
		return func() {}, err
	}

	// Listen on NATS controller
	if err := natsController.Listen(); err != nil {
		return func() {}, err
	}

	return func() {
		natsController.Close()
	}, nil
}

func run() int {
	// Init health server
	h := health.New()
	port, err := strconv.Atoi(os.Getenv("HEALTH_PORT"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("getting health port: %w", err))
		return 255
	}
	go h.HTTPServe(port)

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	component, err := initApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initialize application: %w", err))
		return 255
	}

	// Init controller
	cleanupController, err := initController(component)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initializing NATS controller: %w", err))
		return 255
	}
	defer cleanupController()

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
