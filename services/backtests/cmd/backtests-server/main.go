package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers/http/health"
	"github.com/digital-feather/cryptellation/services/backtests/internal/infrastructure/db/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/infrastructure/pubsub/nats"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
)

func initApp() (*application.Application, func(), error) {
	// Init database client
	db, err := redis.New()
	if err != nil {
		return nil, func() {}, err
	}

	// Init pubsub client
	ps, err := nats.New()
	if err != nil {
		return nil, func() {}, err
	}

	// Init candlestick client
	cs, closeCsClient, err := client.New()
	if err != nil {
		return nil, func() {}, err
	}
	defer func() {
		// Cleanup candlestick client if there is an error
		if err != nil {
			_ = closeCsClient()
		}
	}()

	// Init application
	app, err := application.New(cs, db, ps)
	if err != nil {
		return nil, func() {}, err
	}

	return app, func() { _ = closeCsClient() }, nil
}

func run() int {
	// Init health server
	h := health.New()
	go h.Serve()

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	app, cleanupApp, err := initApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured when %+v\n", fmt.Errorf("initialize application: %w", err))
		return 255
	}
	defer cleanupApp()

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

	log.Print("service is shutting down...")
	return 0
}

func main() {
	os.Exit(run())
}
