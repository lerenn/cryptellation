module cryptellation/client

go 1.22.4

replace (
	cryptellation/client => ../../clients/go
	cryptellation/internal => ../../internal
	cryptellation/pkg => ../../pkg

	cryptellation/svc/backtests => ../../svc/backtests
	cryptellation/svc/candlesticks => ../../svc/candlesticks
	cryptellation/svc/exchanges => ../../svc/exchanges
	cryptellation/svc/forwardtests => ../../svc/forwardtests
	cryptellation/svc/indicators => ../../svc/indicators
	cryptellation/svc/ticks => ../../svc/ticks
)

require (
	cryptellation/internal v0.0.0-00010101000000-000000000000
	cryptellation/pkg v0.0.0-00010101000000-000000000000
	cryptellation/svc/backtests v0.0.0-00010101000000-000000000000
	cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000
	cryptellation/svc/forwardtests v0.0.0-00010101000000-000000000000
	cryptellation/svc/indicators v0.0.0-00010101000000-000000000000
	cryptellation/svc/ticks v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
)

require (
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)
