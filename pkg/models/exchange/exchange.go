package exchange

import (
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

// Equals checks if two exchanges are equal.
func (e Exchange) Equals(e2 Exchange) bool {
	less := func(x, y any) bool {
		switch xv := x.(type) {
		case string:
			yv := y.(string)
			return xv < yv
		case float64:
			yv := y.(float64)
			return xv < yv
		case time.Time:
			yv := y.(time.Time)
			return xv.Before(yv)
		default:
			return false
		}
	}
	diff := cmp.Diff(e, e2, cmpopts.SortSlices(less))
	return diff == ""
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

// String returns the string representation of the exchange.
func (e Exchange) String() string {
	return fmt.Sprintf("Exchange{Name: %s, Periods: %v, Pairs: %v, Fees: %f, LastSyncTime: %s}",
		e.Name, e.Periods, e.Pairs, e.Fees, e.LastSyncTime)
}

// IsOutdated checks if the exchange is outdated.
func (e Exchange) IsOutdated(expirationDuration time.Duration) bool {
	return time.Now().After(e.LastSyncTime.Add(expirationDuration))
}
