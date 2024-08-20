package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"

	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
)

var (
	ErrNoActualPrice = fmt.Errorf("no actual price")
)

const (
	BalanceSymbol = "USDT"
)

func (f ForwardTests) GetStatus(ctx context.Context, forwardTest uuid.UUID) (forwardtest.Status, error) {
	// Get forward test from db
	ft, err := f.db.ReadForwardTest(ctx, forwardTest)
	if err != nil {
		return forwardtest.Status{}, err
	}

	// Get value for each symbol in accounts
	total := 0.0
	for exchange, account := range ft.Accounts {
		for symbol, balance := range account.Balances {
			if symbol == BalanceSymbol {
				total += balance
				continue
			}

			// Get price
			p := symbol + "-" + BalanceSymbol
			cs, err := f.candlesticks.Read(ctx, client.ReadCandlesticksPayload{
				Exchange: exchange,
				Pair:     p,
				Period:   period.M1,
				Start:    utils.ToReference(time.Now().Add(-time.Minute * 10)),
				End:      utils.ToReference(time.Now()),
				Limit:    1,
			})
			if err != nil {
				return forwardtest.Status{}, err
			}

			_, c, ok := cs.Last()
			if !ok {
				return forwardtest.Status{}, fmt.Errorf("%w: %s", ErrNoActualPrice, p)
			}

			// Calculate value
			total += balance * c.Close
		}
	}

	return forwardtest.Status{
		Balance: total,
	}, nil
}
