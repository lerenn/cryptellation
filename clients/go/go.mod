module github.com/lerenn/cryptellation/clients/go

go 1.21.3

replace github.com/lerenn/cryptellation/pkg => ../../pkg

replace github.com/lerenn/cryptellation/svc/backtests => ../../svc/backtests

replace github.com/lerenn/cryptellation/svc/candlesticks => ../../svc/candlesticks

replace github.com/lerenn/cryptellation/svc/exchanges => ../../svc/exchanges

replace github.com/lerenn/cryptellation/svc/indicators => ../../svc/indicators

replace github.com/lerenn/cryptellation/svc/ticks => ../../svc/ticks

require (
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/backtests v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/indicators v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/ticks v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/lerenn/asyncapi-codegen v0.30.0 // indirect
	github.com/nats-io/nats.go v1.31.0 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	gorm.io/gorm v1.25.5 // indirect
)
