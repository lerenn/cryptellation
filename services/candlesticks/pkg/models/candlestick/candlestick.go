package candlestick

import (
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client/proto"
)

type Candlestick struct {
	Open       float64 `bson:"open"     json:"open,omitempty"`
	High       float64 `bson:"high"     json:"high,omitempty"`
	Low        float64 `bson:"low"      json:"low,omitempty"`
	Close      float64 `bson:"close"    json:"close,omitempty"`
	Volume     float64 `bson:"volume"   json:"volume,omitempty"`
	Uncomplete bool    `bson:"complete" json:"uncomplete,omitempty"`
}

func (cs Candlestick) Equal(b Candlestick) bool {
	o := cs.Open == b.Open
	h := cs.High == b.High
	l := cs.Low == b.Low
	c := cs.Close == b.Close
	v := cs.Volume == b.Volume
	u := cs.Uncomplete == b.Uncomplete
	return o && h && l && c && v && u
}

func FromProtoBuf(pbc *proto.Candlestick) (time.Time, Candlestick, error) {
	t, err := time.Parse(time.RFC3339, pbc.Time)
	if err != nil {
		return time.Time{}, Candlestick{}, err
	}

	return t, Candlestick{
		Open:   float64(pbc.Open),
		High:   float64(pbc.High),
		Low:    float64(pbc.Low),
		Close:  float64(pbc.Close),
		Volume: float64(pbc.Volume),
	}, nil
}

func (cs Candlestick) ToProfoBuff(t time.Time) *proto.Candlestick {
	return &proto.Candlestick{
		Time:   t.UTC().Format(time.RFC3339Nano),
		Open:   float32(cs.Open),
		High:   float32(cs.High),
		Low:    float32(cs.Low),
		Close:  float32(cs.Close),
		Volume: float32(cs.Volume),
	}
}

func (cs Candlestick) PriceByType(pt PriceType) float64 {
	return PriceByType(cs.Open, cs.High, cs.Low, cs.Close, pt)
}
