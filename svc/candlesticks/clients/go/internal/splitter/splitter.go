package splitter

import (
	"context"
	"time"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

type Splitter struct {
	client client.Client
}

func New(client client.Client) client.Client {
	return &Splitter{
		client: client,
	}
}

const (
	ReadPeriods = 50
)

func (s Splitter) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// If start and end are set, check if the request can be done in one request
	if payload.Start != nil && payload.End != nil {
		count := payload.Period.CountBetweenTimes(*payload.Start, *payload.End)
		if count <= ReadPeriods || (payload.Limit > 0 && payload.Limit < ReadPeriods) {
			return s.client.Read(ctx, payload)
		}
	}

	// Otherwise, split the request
	return s.readWithSplit(ctx, payload)
}

func (s Splitter) readWithSplit(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// Set end if not set
	end := time.Now()
	if payload.End != nil {
		end = *payload.End
	}

	// Set start if not set
	if payload.Start == nil {
		payload.Start = utils.ToReference(end.Add(-time.Duration(ReadPeriods) * payload.Period.Duration()))
	}
	payload.End = payload.Start

	// Loop to gather all candlesticks
	finalList := candlestick.NewList(payload.Exchange, payload.Pair, payload.Period)
	for payload.End.Before(end) {
		// Set times
		payload.Start = payload.End
		payload.End = utils.ToReference(payload.End.Add(payload.Period.Duration() * ReadPeriods))
		if end.Before(*payload.End) {
			payload.End = utils.ToReference(end)
		}

		// Read on one call
		m, err := s.client.Read(ctx, payload)
		if err != nil {
			return nil, err
		}

		// Merge lists
		if err := finalList.Merge(m, nil); err != nil {
			return nil, err
		}
	}

	return finalList, nil
}

func (s Splitter) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	return s.client.ServiceInfo(ctx)
}

func (s Splitter) Close(ctx context.Context) {
	s.client.Close(ctx)
}
