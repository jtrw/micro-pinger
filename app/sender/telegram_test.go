package sender

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSenderTelegram(t *testing.T) {
	testCases := []struct {
		name          string
		senderType    string
		message       Message
		expectedError bool
		sendFunc      func() error
	}{
		{
			name:       "ValidSenderType",
			senderType: "telegram",
			message: Message{
				ServiceName: "example",
			},
			expectedError: false,
			sendFunc: func() error {
				return nil
			},
		},
		{
			name:       "UnsupportedSenderType",
			senderType: "unsupported",
			message: Message{
				ServiceName: "example",
			},
			expectedError: true,
			sendFunc: func() error {
				return errors.New("mock error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sender := &MockSender{sendFunc: tc.sendFunc}
			err := sender.Send()

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
