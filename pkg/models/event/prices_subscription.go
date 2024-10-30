package event

type PricesSubscription struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}
