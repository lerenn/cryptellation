package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
)

const (
	HeathPortEnvVar = "HEALTH_PORT"
)

type Health struct {
	isReady atomic.Value
	port    int

	livenessCounter  telemetry.Counter
	readinessCounter telemetry.Counter
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

	desc := "how many times liveness function has been called."
	h.livenessCounter, _ = telemetry.CI("health", "liveness_calls", desc)
	h.readinessCounter, _ = telemetry.CI("health", "readiness_calls", desc)

	return &h, nil
}

func (h *Health) Ready(isReady bool) {
	h.isReady.Store(isReady)
}

func (h *Health) HTTPServe(ctx context.Context) {
	http.HandleFunc("/liveness", h.liveness())
	http.HandleFunc("/readiness", h.readiness())

	url := fmt.Sprintf(":%d", h.port)
	telemetry.L(ctx).Info("Starting: HTTP Health Listener")
	telemetry.L(ctx).Error(http.ListenAndServe(url, nil).Error())
}

func (h *Health) liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add trace
		ctx, tracer := telemetry.T(r.Context(), "health", "liveness")
		defer tracer.End()

		// Add call counter
		defer h.livenessCounter.Add(ctx, 1)

		// Add log
		defer telemetry.L(ctx).Debug("liveness called")

		// Write response
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Health) readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add trace
		ctx, tracer := telemetry.T(r.Context(), "health", "readiness")
		defer tracer.End()

		// Add call counter
		defer h.readinessCounter.Add(ctx, 1)

		// Add log
		defer telemetry.L(ctx).Debug("readiness called")

		// Write response
		if !h.isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
