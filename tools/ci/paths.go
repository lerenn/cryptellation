package main

var (
	pathCmdCli = "/cmd/cryptellation"
	pathCmdTui = "/cmd/cryptellation-tui"

	pathPkg = "/pkg"

	pathSvcBacktests    = "/svc/backtests"
	pathSvcCandlesticks = "/svc/candlesticks"
	pathSvcExchanges    = "/svc/exchanges"
	pathSvcIndicators   = "/svc/indicators"
	pathSvcTicks        = "/svc/ticks"

	pathToolsCi = "/tools/ci"
)

var (
	pathServices = []string{
		pathSvcBacktests,
		pathSvcCandlesticks,
		pathSvcExchanges,
		pathSvcIndicators,
		pathSvcTicks,
	}
)

var (
	pathModules = []string{
		pathCmdCli,
		pathCmdTui,

		pathPkg,

		pathSvcBacktests,
		pathSvcCandlesticks,
		pathSvcExchanges,
		pathSvcIndicators,
		pathSvcTicks,

		pathToolsCi,
	}
)
