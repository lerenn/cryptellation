package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/lerenn/cryptellation/internal/core/candlesticks"
	"github.com/lerenn/cryptellation/internal/ctrl/candlesticks/events"
	"github.com/lerenn/cryptellation/internal/infra/candlesticks/db/sql"
	"github.com/lerenn/cryptellation/internal/infra/candlesticks/exchanges"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/http/health"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Launch the service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve()
	},
}

func initApp() (candlesticks.Interface, error) {
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
	return candlesticks.New(db, exchanges), nil
}

func initController(component candlesticks.Interface) (func(), error) {
	// Init NATS controller
	natsController, err := events.NewNATS(config.LoadNATSConfigFromEnv(), component)
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

func serve() error {
	// Init health server
	h := health.New()
	port, err := strconv.Atoi(os.Getenv("HEALTH_PORT"))
	if err != nil {
		return fmt.Errorf("getting health port: %w", err)
	}
	go h.HTTPServe(port)

	// Listen interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Init application
	component, err := initApp()
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	// Init controller
	cleanupController, err := initController(component)
	if err != nil {
		return fmt.Errorf("initializing NATS controller: %w", err)
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
	return nil
}
