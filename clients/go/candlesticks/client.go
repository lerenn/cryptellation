// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=client.go -destination=client.mock.gen.go -package candlesticks

package candlesticks

type ClientConfig struct {
}

type Client interface {
	// ReadCandlesticks(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error)
}

type ReadCandlesticksPayload struct {
	// candlesticks.GetCachedPayload
}
