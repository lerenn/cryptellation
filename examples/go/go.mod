module github.com/lerenn/cryptellation/examples/go

go 1.22.4

replace (
	github.com/lerenn/cryptellation/clients/go => ../../clients/go
	github.com/lerenn/cryptellation/internal => ../../internal
	github.com/lerenn/cryptellation/pkg => ../../pkg

	github.com/lerenn/cryptellation/client => ../../svc/backtests
	github.com/lerenn/cryptellation/svc/candlesticks => ../../svc/candlesticks
	github.com/lerenn/cryptellation/svc/exchanges => ../../svc/exchanges
	github.com/lerenn/cryptellation/svc/forwardtests => ../../svc/forwardtests
	github.com/lerenn/cryptellation/svc/indicators => ../../svc/indicators
	github.com/lerenn/cryptellation/svc/ticks => ../../svc/ticks
)

require (
	github.com/lerenn/cryptellation/clients/go v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/internal v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/client v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/forwardtests v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/indicators v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/ticks v0.0.0-00010101000000-000000000000
)

require (
	github.com/lerenn/cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000 // indirect
	github.com/agoda-com/opentelemetry-logs-go v0.5.1 // indirect
	github.com/agoda-com/otelzap v0.1.1 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/go-echarts/go-echarts/v2 v2.4.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
