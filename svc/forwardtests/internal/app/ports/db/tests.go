package db

import (
	"github.com/stretchr/testify/suite"
)

type ForwardTestSuite struct {
	suite.Suite
	DB Port
}
