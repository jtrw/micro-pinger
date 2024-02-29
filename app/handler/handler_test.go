package handler

import (
	//"log"
	"io/ioutil"
	config "micro-pinger/v2/app/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	StatusCode int
	Body       string
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Customize the response based on your test scenario
	// For example, simulate a successful response
	return &http.Response{
		StatusCode: m.StatusCode, //http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(m.Body)),
	}, nil
}

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

func TestCheckService(t *testing.T) {
	serviceName := "SampleService_2"
	sampleService := config.Service{
		Name: serviceName,
		URL:  "https://example.com",
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

	for i := 0; i < 3; i++ {
		mockClient := &MockHTTPClient{StatusCode: http.StatusBadRequest, Body: "Bad Request"}
		handler.CheckService(mockClient, sampleService)
	}

	assert.Equal(t, 3, FailureThreshold[serviceName+"_SampleAlert"])
	assert.Equal(t, 0, SuccessThreshold[serviceName+"_SampleAlert"])

	mockClient := &MockHTTPClient{StatusCode: http.StatusOK, Body: "OK"}
	handler.CheckService(mockClient, sampleService)

	assert.Equal(t, 3, FailureThreshold[serviceName+"_SampleAlert"])
	assert.Equal(t, 1, SuccessThreshold[serviceName+"_SampleAlert"])

	mockClient = &MockHTTPClient{StatusCode: http.StatusOK, Body: "OK"}
	handler.CheckService(mockClient, sampleService)

	assert.Equal(t, 0, FailureThreshold[serviceName+"_SampleAlert"])
	assert.Equal(t, 0, SuccessThreshold[serviceName+"_SampleAlert"])
}
