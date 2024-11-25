package exchange

import (
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/utils"
)

// Exchange represents an exchange.
type Exchange struct {
	Name         string
	Periods      []string
	Pairs        []string
	Fees         float64
	LastSyncTime time.Time
}

// Merge merges two exchanges into one.
func (e Exchange) Merge(e2 Exchange) Exchange {
	return Exchange{
		Name:    e.Name,
		Periods: utils.MergeSliceIntoUnique(e.Periods, e2.Periods),
		Pairs:   utils.MergeSliceIntoUnique(e.Pairs, e2.Pairs),
		Fees:    e.Fees,
	}
}

// AddPair adds a pair to the exchange.
func (e *Exchange) AddPair(symbols ...string) {
	e.Pairs = utils.MergeSliceIntoUnique(e.Pairs, symbols)
}

// AddPeriods adds a period to the exchange.
func (e *Exchange) AddPeriods(symbols ...string) {
	e.Periods = utils.MergeSliceIntoUnique(e.Periods, symbols)
}

// ArrayToMap converts an array of exchanges to a map of exchanges.
func ArrayToMap(exchanges []Exchange) map[string]Exchange {
	mappedExchanges := make(map[string]Exchange, len(exchanges))
	for _, exch := range exchanges {
		mappedExchanges[exch.Name] = exch
	}
	return mappedExchanges
}

// MapToArray converts a map of exchanges to an array of exchanges.
func MapToArray(mappedExchanges map[string]Exchange) []Exchange {
	exchanges := make([]Exchange, 0, len(mappedExchanges))
	for _, exch := range mappedExchanges {
		exchanges = append(exchanges, exch)
	}
	return exchanges
}

// GetExpiredExchangesNames returns the names of the exchanges that are expired.
func GetExpiredExchangesNames(
	expectedExchanges []string,
	exchangesFromDB []Exchange,
	expirationDuration time.Duration,
) (toSync []string, err error) {
	mappedExchanges := ArrayToMap(exchangesFromDB)

	toSync = make([]string, 0, len(expectedExchanges))
	for _, name := range expectedExchanges {
		exch, ok := mappedExchanges[name]
		if ok && time.Now().Before(exch.LastSyncTime.Add(expirationDuration)) {
			continue
		}

		toSync = append(toSync, name)
	}

	return toSync, err
}
