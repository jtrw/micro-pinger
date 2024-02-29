package handler

import (
	//"log"
	config "micro-pinger/v2/app/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Check(t *testing.T) {
	// Create a sample handler with one service for testing
	sampleService := config.Service{
		Name: "SampleService",
		URL:  "https://bad123456.com",
		Response: config.Response{
			Status: http.StatusOK,
			Body:   "OK",
		},
		Alerts: []config.Alert{
			{
				Name:          "SampleAlert",
				Webhook:       "https://hooks.slack.com/services/123456/7890",
				Type:          "slack",
				Failure:       3,
				Success:       2,
				SendOnResolve: true,
			},
		},
	}
	handler := NewHandler([]config.Service{sampleService})

	// Create a mock server to simulate successful responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer mockServer.Close()

	// Set the mock server URL as the service URL
	sampleService.URL = mockServer.URL

	// Perform the Check operation
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/check", nil)
		handler.Check(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	}

	time.Sleep(1 * time.Second)

	if FailureThreshold["SampleService_SampleAlert"] != 3 {
		t.Errorf("Expected FailureThreshold to be 3, got %d", FailureThreshold["SampleService_SampleAlert"])
	}
}
