module github.com/lerenn/cryptellation/cmd/cryptellation

go 1.21.3

replace github.com/lerenn/cryptellation/clients/go => ../../clients/go

replace github.com/lerenn/cryptellation/pkg => ../../pkg

replace github.com/lerenn/cryptellation/svc/backtests => ../../svc/backtests

replace github.com/lerenn/cryptellation/svc/candlesticks => ../../svc/candlesticks

replace github.com/lerenn/cryptellation/svc/exchanges => ../../svc/exchanges

replace github.com/lerenn/cryptellation/svc/indicators => ../../svc/indicators

replace github.com/lerenn/cryptellation/svc/ticks => ../../svc/ticks

require (
	github.com/lerenn/cryptellation/clients/go v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/ticks v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-git/go-git/v5 v5.11.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/lerenn/asyncapi-codegen v0.30.2 // indirect
	github.com/lerenn/cryptellation/svc/backtests v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000 // indirect
	github.com/lerenn/cryptellation/svc/indicators v0.0.0-00010101000000-000000000000 // indirect
	github.com/nats-io/nats.go v1.31.0 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/skeema/knownhosts v1.2.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gorm.io/gorm v1.25.5 // indirect
)
