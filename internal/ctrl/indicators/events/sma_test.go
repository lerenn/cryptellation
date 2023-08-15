package events

import (
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/stretchr/testify/suite"
)

func TestSMASuite(t *testing.T) {
	suite.Run(t, new(SMASuite))
}

type SMASuite struct {
	suite.Suite
}

func (suite *SMASuite) TestSetSmaResponseMessage() {
	resp := NewSmaResponseMessage()
	resp.Set(timeserie.New[float64]().
		Set(time.Unix(0, 0), 1).
		Set(time.Unix(60, 0), 2))

	suite.Require().NotNil(resp.Payload.Data)
	suite.Require().Len(*resp.Payload.Data, 2)
}
