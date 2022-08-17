package status

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}

type StatusSuite struct {
	suite.Suite
}

func (suite *StatusSuite) TestFromJSON() {
	tests := []struct {
		Input       []byte
		Output      Status
		OutputError bool
	}{
		{
			Input:       []byte("{\"finished\":true}"),
			Output:      Status{Finished: true},
			OutputError: false,
		},
		{
			Input:       []byte("{\"finished\":false}"),
			Output:      Status{Finished: false},
			OutputError: false,
		},
		{
			Input:       []byte("{\"finished\":dfdsqfq}"),
			Output:      Status{},
			OutputError: true,
		},
	}

	for i, t := range tests {
		output, err := FromJSON(t.Input)
		if t.OutputError {
			suite.Require().Error(err, i)
			continue
		} else {
			suite.Require().NoError(err, i)
		}

		suite.Require().Equal(t.Output, output, i)
	}
}

func (suite *StatusSuite) TestMarshalingJSON() {
	as := suite.Require()

	s := Status{
		Finished: true,
	}

	b, err := json.Marshal(s)
	as.NoError(err)

	s2 := Status{}
	as.NoError(json.Unmarshal(b, &s2))
	as.Equal(s, s2)
}
