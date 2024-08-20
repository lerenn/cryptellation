module github.com/lerenn/cryptellation/cmd/cryptellation

go 1.22.4

replace (
	github.com/lerenn/cryptellation/clients/go => ../../clients/go
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
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/lerenn/cryptellation/client v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/forwardtests v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/indicators v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/ticks v0.0.0-00010101000000-000000000000 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/cloudflare/circl v1.3.9 // indirect
	github.com/cyphar/filepath-securejoin v0.3.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-git/go-git/v5 v5.12.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)
