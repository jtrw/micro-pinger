package sender

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockSender є моком для інтерфейсу Sender.
type MockSender struct {
	sendFunc func() error
}

// Send викликає функцію sendFunc моку.
func (m *MockSender) Send() error {
	return m.sendFunc()
}

func TestNewSender(t *testing.T) {
	testCases := []struct {
		name          string
		senderType    string
		message       Message
		expectedError bool
		sendFunc      func() error
	}{
		{
			name:       "ValidSenderType",
			senderType: "slack",
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

func TestSlack_Send(t *testing.T) {
	testCases := []struct {
		name          string
		message       Message
		sendFunc      func() error
		expectedError bool
	}{
		{
			name: "SuccessfulSend",
			message: Message{
				Webhook: "https://example.com/webhook",
			},
			sendFunc: func() error {
				return nil
			},
			expectedError: false,
		},
		{
			name: "SendError",
			message: Message{
				Webhook: "https://example.com/webhook",
			},
			sendFunc: func() error {
				return errors.New("mock error")
			},
			expectedError: true,
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
