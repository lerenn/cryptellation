// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=client.go -destination=mock.gen.go -package temporal

package temporal

import (
	temporalclient "go.temporal.io/sdk/client"
)

// Client is the client interface from temporal library.
type Client interface {
	temporalclient.Client
}
