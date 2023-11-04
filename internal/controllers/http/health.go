package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/agoda-com/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	HeathPortEnvVar = "HEALTH_PORT"
)

type Health struct {
	isReady atomic.Value
	port    int

	livenessCounter  metric.Int64Counter
	readinessCounter metric.Int64Counter
}

func NewHealth() (*Health, error) {
	var h Health

	// Get health port
	port, err := strconv.Atoi(os.Getenv(HeathPortEnvVar))
	if err != nil {
		return nil, fmt.Errorf("getting health port: %w", err)
	}
	h.port = port

	// Readiness to false
	h.isReady.Store(false)

	// Initialize telemetry
	desc := metric.WithDescription("how many times liveness function has been called.")
	h.livenessCounter, _ = otel.Meter("health").Int64Counter("liveness_calls", desc)
	h.readinessCounter, _ = otel.Meter("health").Int64Counter("readiness_calls", desc)

	return &h, nil
}

func (h *Health) Ready(isReady bool) {
	h.isReady.Store(isReady)
}

func (h *Health) HTTPServe() {
	http.HandleFunc("/liveness", h.liveness())
	http.HandleFunc("/readiness", h.readiness())

	url := fmt.Sprintf(":%d", h.port)
	log.Println("Starting: HTTP Health Listener")
	log.Fatal(http.ListenAndServe(url, nil))
}

func (h *Health) liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add trace
		ctx, span := otel.Tracer("health").Start(r.Context(), "liveness")
		defer span.End()

		// Add call counter
		defer h.livenessCounter.Add(ctx, 1)

		// Add log
		defer otelzap.Ctx(ctx).Debug("liveness called")

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Health) readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add trace
		ctx, span := otel.Tracer("health").Start(r.Context(), "readiness")
		defer span.End()

		// Add call counter
		defer h.readinessCounter.Add(ctx, 1)

		// Add log
		defer otelzap.Ctx(ctx).Debug("readiness called")

		// Write response
		if !h.isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
