package sender

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
	// Create a sample message for testing
	message := Message{
		Status:      "OK",
		Webhook:     "https://example.com/slack-webhook",
		Datetime:    "2024-02-28T12:34:56",
		Url:         "https://example.com",
		ServiceName: "TestService",
		Response: Response{
			Text: "Test message",
			Err:  nil,
			Code: 200,
		},
	}

	// Create a mock server to simulate Slack API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and content type
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read request body
		var receivedMessage SlackMessage
		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Compare the received message with the expected message
		expectedMessage := getTextMessage(message)
		if receivedMessage.Text != expectedMessage {
			t.Errorf("Expected message: %s, got: %s", expectedMessage, receivedMessage.Text)
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer mockServer.Close()

	// Set the mock server URL as the Slack webhook
	message.Webhook = mockServer.URL

	// Create a Slack sender with the mock message
	slackSender := NewSlack(message)

	// Call the Send method
	err := slackSender.Send()

	// Check for errors
	if err != nil {
		t.Errorf("Error sending message to Slack: %v", err)
	}
}
