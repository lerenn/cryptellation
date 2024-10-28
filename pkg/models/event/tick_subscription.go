package event

type TickSubscription struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}
