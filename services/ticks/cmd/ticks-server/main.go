package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/services/ticks/internal/application"
	exchPorts "github.com/digital-feather/cryptellation/services/ticks/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/ticks/internal/controllers/http/health"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/exchanges"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/exchanges/binance"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/ticks/internal/infrastructure/vdb/redis"
)

func initApp() (*application.Application, error) {
	// Set exchange services
	binanceService, err := binance.New()
	if err != nil {
		return nil, err
	}

	// Regroup every exchange service
	services := map[string]exchPorts.Adapter{
		exchanges.BinanceName: binanceService,
	}

	db, err := redis.New()
	if err != nil {
		return nil, err
	}

	ps, err := nats.New()
	if err != nil {
		return nil, err
	}

	return application.New(db, ps, services)
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
	grpcController := grpc.New(app)
	if err := grpcController.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("running application: %w", err))
		return 255
	}
	defer grpcController.GracefulStop()

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
