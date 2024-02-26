package handler

import (
	//"bytes"
	//"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	//"io/ioutil"
	config "micro-pinger/v2/app/service"
)

// MockSender is a mock implementation of the Sender interface.
type MockSender struct {
	sendFunc func() error
}

func (m *MockSender) Send() error {
	return m.sendFunc()
}

func TestHandler_Check(t *testing.T) {
	// Test case for Handler.Check function
	t.Run("CheckServices", func(t *testing.T) {
		// Mock HTTP request
		req, err := http.NewRequest("GET", "/check", nil)
		assert.NoError(t, err)

		// Mock HTTP response recorder
		w := httptest.NewRecorder()

		// Create a Handler with services
		services := []config.Service{
			// Add your test services here
		}
		handler := NewHandler(services)

		// Call Check function
		handler.Check(w, req)

		// Check HTTP response status
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the response body (assuming you have a valid JSON response)
		expectedResponse := `{"status":"ok"}` + "\n"
		assert.Equal(t, expectedResponse, w.Body.String())
	})
}
