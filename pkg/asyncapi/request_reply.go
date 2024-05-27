package asyncapi

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// AddReplyToSuffix adds a random suffix to the replyTo channel name.
func AddReplyToSuffix(address string) string {
	return fmt.Sprintf("%s.%s", address, uuid.New().String())
}

func UnwrapError(errAny any) error {
	// Check if the error is not nil
	if errAny == nil {
		return nil
	}

	// Serialize the error
	b, err := json.Marshal(errAny)
	if err != nil {
		return err
	}

	// Deserialize the error
	var errMsg ErrorSchema
	if err := json.Unmarshal(b, &errMsg); err != nil {
		return err
	}

	// Check if error is filled
	if (errMsg.Code == 0 || (errMsg.Code >= 200 && errMsg.Code < 300)) && errMsg.Message == "" {
		return nil
	}

	// Return the error
	return fmt.Errorf("%s (code=%d)", errMsg.Message, errMsg.Code)
}
