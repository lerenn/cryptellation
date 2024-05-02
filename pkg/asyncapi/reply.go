package asyncapi

import (
	"fmt"
	"math/rand/v2"
)

const (
	// RandomSuffixSize is the suffix size used to generate a random replyTo channel
	RandomSuffixSize = 1024
)

// AddReplyToSuffix adds a random suffix to the replyTo channel name.
func AddReplyToSuffix(address string) string {
	return fmt.Sprintf("%s.%d", address, rand.IntN(RandomSuffixSize))
}
