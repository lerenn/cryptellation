package asyncapi

import "github.com/lerenn/cryptellation/forwardtests/pkg/forwardtest"

func (msg *StatusResponseMessage) Set(status forwardtest.Status) {
	msg.Payload.Status = &ForwardTestStatusSchema{
		Balance: status.Balance,
	}
}

func (msg StatusResponseMessage) ToModel() forwardtest.Status {
	return forwardtest.Status{
		Balance: msg.Payload.Status.Balance,
	}
}
