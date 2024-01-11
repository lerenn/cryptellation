module github.com/lerenn/cryptellation/cmd

go 1.21.3

replace github.com/lerenn/cryptellation/pkg => ../pkg

replace github.com/lerenn/cryptellation/svc/backtests => ../svc/backtests
replace github.com/lerenn/cryptellation/svc/candlesticks => ../svc/candlesticks
replace github.com/lerenn/cryptellation/svc/exchanges => ../svc/exchanges
replace github.com/lerenn/cryptellation/svc/indicators => ../svc/indicators
replace github.com/lerenn/cryptellation/svc/ticks => ../svc/ticks

require (
	dagger.io/dagger v0.9.5
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/99designs/gqlgen v0.17.31 // indirect
	github.com/Khan/genqlient v0.6.0 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vektah/gqlparser/v2 v2.5.6 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
