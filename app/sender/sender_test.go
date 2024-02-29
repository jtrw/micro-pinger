package sender

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSender_NotSupported(t *testing.T) {
	message := Message{
		Status:      "OK",
		Webhook:     "https://example.com/unsupported-webhook",
		Datetime:    "2024-02-28T12:34:56",
		Url:         "https://example.com",
		ServiceName: "TestService",
		Response: Response{
			Text: "Test message",
			Err:  nil,
			Code: 200,
		},
	}

	_, err := NewSender("unsupported", message)
	assert.Error(t, err)
}
