package domain

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/timeserie"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/indicators/pkg/sma"
)

func (ind indicators) GetCachedSMA(ctx context.Context, payload app.GetCachedSMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Get cached SMA from DB
	ts, err := ind.db.GetSMA(ctx, db.ReadSMAPayload{
		Exchange:     payload.Exchange,
		Pair:         payload.Pair,
		Period:       payload.Period,
		PeriodNumber: payload.PeriodNumber,
		PriceType:    payload.PriceType,
		Start:        payload.Start,
		End:          payload.End,
	})
	if err != nil {
		return ts, err
	}

	// Check if current candlestick will be requested
	// If that's the case, we'll need to recalculate the SMA as the value has changed
	requested := payload.Period.RoundTime(payload.End)
	roundedNow := payload.Period.RoundTime(time.Now())
	possiblyOutdatedSMA := requested.Equal(roundedNow)

	// Check if there is missing points
	missingPoints := ts.AreMissing(payload.Start, payload.End, payload.Period.Duration(), 0)

	// Check if we can return or if we can return right now
	if !missingPoints && !possiblyOutdatedSMA {
		return ts, nil
	}

	// Generate SMA points
	ts, err = ind.generateSMA(ctx, payload)
	if err != nil {
		return ts, err
	}

	// Save SMA points to DB and return the result
	return ts, ind.db.UpsertSMA(ctx, db.WriteSMAPayload{
		Exchange:     payload.Exchange,
		Pair:         payload.Pair,
		Period:       payload.Period,
		PeriodNumber: payload.PeriodNumber,
		PriceType:    payload.PriceType,
		TimeSerie:    ts,
	})
}

func (ind indicators) generateSMA(
	ctx context.Context,
	payload app.GetCachedSMAPayload,
) (*timeserie.TimeSerie[float64], error) {
	// Get necessary candlesticks
	cs, err := ind.candlesticks.Read(ctx, client.ReadCandlesticksPayload{
		Exchange: payload.Exchange,
		Pair:     payload.Pair,
		Period:   payload.Period,
		Start:    utils.ToReference(payload.Start.Add(-payload.Period.Duration() * time.Duration(payload.PeriodNumber))),
		End:      utils.ToReference(payload.End),
	})
	if err != nil {
		return &timeserie.TimeSerie[float64]{}, err
	}

	// Generate SMAs and return them
	return sma.TimeSerie(sma.TimeSeriePayload{
		Candlesticks: cs,
		PriceType:    payload.PriceType,
		Start:        payload.Start,
		End:          payload.End,
		PeriodNumber: payload.PeriodNumber,
	}), nil
}
