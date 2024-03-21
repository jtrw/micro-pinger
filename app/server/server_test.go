package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_Run(t *testing.T) {
	t.Run("ServerRun", func(t *testing.T) {
		testServer := Server{
			Listen: "localhost:8080",
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := testServer.Run(ctx)
			assert.Equal(t, http.ErrServerClosed, err, "Expected http.ErrServerClosed")
		}()
	})
}

func TestServer_Routes(t *testing.T) {
	t.Run("ServerRoutes", func(t *testing.T) {
		testServer := Server{}

		req, err := http.NewRequest("GET", "/api/v1/check", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		router := testServer.routes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	})

	t.Run("ServerRoutesPing", func(t *testing.T) {
		testServer := Server{}

		req, err := http.NewRequest("GET", "/ping", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		router := testServer.routes()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
		assert.Equal(t, "pong", w.Body.String(), "Expected response body 'pong'")
	})
}

func TestRest_RobotsCheck(t *testing.T) {
	srv := Server{Listen: "localhost:54009", Version: "v1", Secret: "12345"}

	ts := httptest.NewServer(srv.routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/robots.txt")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "User-agent: *\nDisallow: /\n", string(body))
}
