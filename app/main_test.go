// main_test.go
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	port := 40000 + int(rand.Int31n(10000))
	os.Args = []string{"app", "--secret=123", "--config=../config.default.yaml", "--listen=" + "localhost:" + strconv.Itoa(port)}

	done := make(chan struct{})
	go func() {
		<-done
		e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.NoError(t, e)
	}()

	finished := make(chan struct{})
	go func() {
		main()
		close(finished)
	}()

	defer func() {
		close(done)
		<-finished
	}()

	waitForHTTPServerStart(port)
	time.Sleep(time.Second)
	client := &http.Client{}

	{
		url := fmt.Sprintf("http://localhost:%d/ping", port)
		req, err := getRequest(url)
		req.Header.Set("Api-Key", "123")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "pong", string(body))
	}
}

func waitForHTTPServerStart(port int) {
	client := http.Client{Timeout: time.Second}
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		if resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ping", port)); err == nil {
			_ = resp.Body.Close()
			return
		}
	}
}

func getRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	return req, err
}
