package candlestick

import "fmt"

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

func (cs Candlestick) PriceByType(pt PriceType) float64 {
	return PriceByType(cs.Open, cs.High, cs.Low, cs.Close, pt)
}

func (cs Candlestick) String() string {
	return fmt.Sprintf("O:%f | H:%f | L:%f | C:%f", cs.Open, cs.High, cs.Low, cs.Close)
}
