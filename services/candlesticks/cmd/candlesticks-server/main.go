package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/http/health"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/nats"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/infrastructure/db/sql"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/infrastructure/exchanges/binance"
)

func initApp() (*application.Application, error) {
	// Init database client
	db, err := sql.New()
	if err != nil {
		return nil, err
	}

	// Init exchanges connections
	binanceService, err := binance.New()
	if err != nil {
		return nil, err
	}

	// Assembling all services in a map
	services := map[string]exchanges.Adapter{
		binance.Name: binanceService,
	}

	// Init application
	app, err := application.New(db, services)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	app, err := initApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initializing application: %w", err))
		return 255
	}

	// Init grpc server
	natsController := nats.New(app)
	if err := natsController.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running nats controller: %w", err))
		return 255
	}
	defer natsController.Close()

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
