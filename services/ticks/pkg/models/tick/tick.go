package tick

import (
	"encoding/json"
	"time"

	"github.com/digital-feather/cryptellation/services/ticks/pkg/client/proto"
)

type Tick struct {
	Time       time.Time `json:"time"`
	PairSymbol string    `json:"pair_symbol"`
	Price      float32   `json:"price"`
	Exchange   string    `json:"exchange"`
}

func (t Tick) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tick) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t Tick) ToProtoBuf() *proto.Tick {
	return &proto.Tick{
		Time:       t.Time.Format(time.RFC3339Nano),
		Exchange:   t.Exchange,
		PairSymbol: t.PairSymbol,
		Price:      float32(t.Price),
	}
}

func FromProtoBuf(pb *proto.Tick) (Tick, error) {
	t, err := time.Parse(time.RFC3339Nano, pb.Time)
	if err != nil {
		return Tick{}, err
	}

	return Tick{
		Time:       t,
		Exchange:   pb.Exchange,
		PairSymbol: pb.PairSymbol,
		Price:      pb.Price,
	}, nil
}
