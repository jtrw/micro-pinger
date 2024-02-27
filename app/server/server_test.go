package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Run(t *testing.T) {
	// Test case for Server.Run function
	t.Run("ServerRun", func(t *testing.T) {
		// Create a Server instance for testing
		testServer := Server{
			Listen: "localhost:8080",
			// Add other necessary fields
		}

		// Create a context with a cancellation function
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Run the server in a goroutine
		go func() {
			err := testServer.Run(ctx)
			assert.Equal(t, http.ErrServerClosed, err, "Expected http.ErrServerClosed")
		}()

		// Make a request to the server or perform other actions as needed
		// For example, you can use an HTTP client to send requests and verify responses.

		// To stop the server, cancel the context (e.g., cancel())

		// Add assertions based on your expectations
	})

	// Add more test cases as needed
}

func TestServer_Routes(t *testing.T) {
	// Test case for Server.routes function
	t.Run("ServerRoutes", func(t *testing.T) {
		// Create a Server instance for testing
		testServer := Server{
			// Add other necessary fields
		}

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/api/v1/check", nil)
		assert.NoError(t, err)

		// Create a test HTTP response recorder
		w := httptest.NewRecorder()

		// Get the router and serve the request
		router := testServer.routes()
		router.ServeHTTP(w, req)

		// Perform assertions based on your expectations
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
		// Add more assertions as needed
	})

	t.Run("ServerRoutesPing", func(t *testing.T) {
		// Create a Server instance for testing
		testServer := Server{
			// Add other necessary fields
		}

		// Create a test HTTP request
		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)

		// Create a test HTTP response recorder
		w := httptest.NewRecorder()

		// Get the router and serve the request
		router := testServer.routes()
		router.ServeHTTP(w, req)

		// Perform assertions based on your expectations
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
		assert.Equal(t, "pong", w.Body.String(), "Expected response body 'pong'")
	})
}
