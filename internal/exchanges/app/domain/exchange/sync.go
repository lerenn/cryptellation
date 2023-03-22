package exchange

import (
	"time"

	"github.com/digital-feather/cryptellation/pkg/types/exchange"
)

const DefaultExpirationDuration = time.Hour

func GetExpiredExchangesNames(
	expectedExchanges []string,
	exchangesFromDB []exchange.Exchange,
) (toSync []string, err error) {
	mappedExchanges := exchange.ArrayToMap(exchangesFromDB)

	toSync = make([]string, 0, len(expectedExchanges))
	for _, name := range expectedExchanges {
		exch, ok := mappedExchanges[name]
		if ok && time.Now().Before(exch.LastSyncTime.Add(DefaultExpirationDuration)) {
			continue
		}

		toSync = append(toSync, name)
	}

	return toSync, err
}
