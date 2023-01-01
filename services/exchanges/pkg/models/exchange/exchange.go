package exchange

import (
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/clients/go/proto"
)

// swagger:model Exchange
type Exchange struct {
	Name           string
	PeriodsSymbols []string
	PairsSymbols   []string
	Fees           float64
	LastSyncTime   time.Time
}

func (e Exchange) Merge(e2 Exchange) Exchange {
	return Exchange{
		Name:           e.Name,
		PeriodsSymbols: addtoUniqueArray(e.PeriodsSymbols, e2.PeriodsSymbols...),
		PairsSymbols:   addtoUniqueArray(e.PairsSymbols, e2.PairsSymbols...),
		Fees:           e.Fees,
	}
}

func (e *Exchange) AddPair(symbols ...string) {
	e.PairsSymbols = addtoUniqueArray(e.PairsSymbols, symbols...)
}

func (e *Exchange) AddPeriods(symbols ...string) {
	e.PeriodsSymbols = addtoUniqueArray(e.PeriodsSymbols, symbols...)
}

func addtoUniqueArray(a1 []string, a2 ...string) []string {
	tmp := make([]string, len(a1))
	copy(tmp, a1)

	for _, s2 := range a2 {
		present := false
		for _, s1 := range a1 {
			if s1 == s2 {
				present = true
				break
			}
		}

		if !present {
			tmp = append(tmp, s2)
		}
	}

	return tmp
}

func ArrayToMap(exchanges []Exchange) map[string]Exchange {
	mappedExchanges := make(map[string]Exchange, len(exchanges))
	for _, exch := range exchanges {
		mappedExchanges[exch.Name] = exch
	}
	return mappedExchanges
}

func MapToArray(mappedExchanges map[string]Exchange) []Exchange {
	exchanges := make([]Exchange, 0, len(mappedExchanges))
	for _, exch := range mappedExchanges {
		exchanges = append(exchanges, exch)
	}
	return exchanges
}

func FromProtoBuf(pb *proto.Exchange) (Exchange, error) {
	lastSyncTime, err := time.Parse(time.RFC3339Nano, pb.LastSyncTime)
	if err != nil {
		return Exchange{}, err
	}

	return Exchange{
		Name:           pb.Name,
		PeriodsSymbols: pb.Periods,
		PairsSymbols:   pb.Pairs,
		Fees:           float64(pb.Fees),
		LastSyncTime:   lastSyncTime,
	}, nil
}

func (e Exchange) ToProfoBuff() *proto.Exchange {
	return &proto.Exchange{
		Name:         e.Name,
		Periods:      e.PeriodsSymbols,
		Pairs:        e.PairsSymbols,
		Fees:         float64(e.Fees),
		LastSyncTime: e.LastSyncTime.UTC().Format(time.RFC3339Nano),
	}
}
