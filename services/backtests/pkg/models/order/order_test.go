package order

import (
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/clients/go/proto"
	"github.com/stretchr/testify/suite"
)

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

type OrderTestSuite struct {
	suite.Suite
}

func (suite *OrderTestSuite) TestFromProtoBuf() {
	pb := &proto.Order{
		Id:            1,
		ExecutionTime: OptString("1970-01-01T00:01:00Z"),
		Type:          TypeIsMarket.String(),
		ExchangeName:  "exchange",
		PairSymbol:    "ETC-USDT",
		Side:          SideIsBuy.String(),
		Quantity:      2,
		Price:         3.0,
	}

	o, err := FromProtoBuf(pb)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(1), o.ID)
	suite.Require().WithinDuration(time.Unix(60, 0), *o.ExecutionTime, time.Second)
	suite.Require().Equal(TypeIsMarket, o.Type)
	suite.Require().Equal("exchange", o.ExchangeName)
	suite.Require().Equal("ETC-USDT", o.PairSymbol)
	suite.Require().Equal(SideIsBuy, o.Side)
	suite.Require().Equal(2.0, o.Quantity)
	suite.Require().Equal(3.0, o.Price)
}

func (suite *OrderTestSuite) TestToProtoBuf() {
	o := Order{
		ID:            1,
		ExecutionTime: OptTime(time.Unix(60, 0)),
		Type:          TypeIsMarket,
		ExchangeName:  "exchange",
		PairSymbol:    "ETC-USDT",
		Side:          SideIsBuy,
		Quantity:      2,
		Price:         3.0,
	}

	pb := o.ToProtoBuf()
	suite.Require().Equal(uint64(1), pb.Id)
	suite.Require().Equal("1970-01-01T00:01:00Z", *pb.ExecutionTime)
	suite.Require().Equal(TypeIsMarket.String(), pb.Type)
	suite.Require().Equal("exchange", pb.ExchangeName)
	suite.Require().Equal("ETC-USDT", pb.PairSymbol)
	suite.Require().Equal(SideIsBuy.String(), pb.Side)
	suite.Require().Equal(float64(2.0), pb.Quantity)
	suite.Require().Equal(float64(3.0), pb.Price)
}

func OptTime(t time.Time) *time.Time {
	return &t
}

func OptString(s string) *string {
	return &s
}
