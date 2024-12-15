package tick

// Subscription is the struct for a tick subscription.
type Subscription struct {
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
}
