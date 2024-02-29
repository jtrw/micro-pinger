// main_test.go
package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"
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

	// defer cleanup because require check below can fail
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
	// wait for up to 10 seconds for server to start before returning it
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
