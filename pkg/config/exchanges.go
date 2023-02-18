package config

type Exchanges struct {
	Binance Binance
}

func LoadExchangesConfigFromEnv() (c Exchanges) {
	c.Binance = LoadBinanceConfigFromEnv()
	return c
}
