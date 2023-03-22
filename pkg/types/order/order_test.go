package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

type OrderTestSuite struct {
	suite.Suite
}
