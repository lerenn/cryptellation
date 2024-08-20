module github.com/lerenn/cryptellation/clients/go

go 1.22.4

replace (
	github.com/lerenn/cryptellation/candlesticks => ../../svc/candlesticks

	github.com/lerenn/cryptellation/client => ../../svc/backtests
	github.com/lerenn/cryptellation/clients/go => ../../clients/go
	github.com/lerenn/cryptellation/exchanges => ../../svc/exchanges
	github.com/lerenn/cryptellation/forwardtests => ../../svc/forwardtests
	github.com/lerenn/cryptellation/indicators => ../../svc/indicators
	github.com/lerenn/cryptellation/internal => ../../internal
	github.com/lerenn/cryptellation/pkg => ../../pkg
	github.com/lerenn/cryptellation/ticks => ../../svc/ticks
)

require (
	github.com/google/uuid v1.6.0
	github.com/lerenn/cryptellation/candlesticks v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/client v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/exchanges v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/forwardtests v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/indicators v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/ticks v0.0.0-00010101000000-000000000000
)

require (
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/lerenn/cryptellation/internal v0.0.0-00010101000000-000000000000 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)
