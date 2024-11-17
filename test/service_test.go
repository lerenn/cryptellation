package test

import "context"

func (suite *EndToEndSuite) TestServiceInfoWorkflow() {
	// Call service info
	info, err := suite.client.Info(context.Background())
	suite.Require().NoError(err)
	suite.Require().NotNil(info)
}
