// main_test.go
package main

import (
	"context"
	"testing"
	//	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	server "micro-pinger/v2/app/server"
	config "micro-pinger/v2/app/service"
)

type MockConfigLoader struct {
	mock.Mock
}

func (m *MockConfigLoader) LoadConfig(filename string) (config.Config, error) {
	args := m.Called(filename)
	return args.Get(0).(config.Config), args.Error(1)
}

func TestMain(t *testing.T) {
	// Mock ConfigLoaderFunc
	mockConfigLoader := new(MockConfigLoader)
	//loadConfigFunc = mockConfigLoader.LoadConfig

	// Expect LoadConfig to be called with "test-config.yml" and return a mock configuration
	mockConfig := config.Config{ /* mocked configuration data */ }
	mockConfigLoader.On("LoadConfig", "config.default.yaml").Return(mockConfig, nil)

	// Run the main function in a goroutine
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()

		// Set the config file to "test-config.yml" for the test
		opts := Options{Config: "config.default.yaml"}
		err := mainWithOptions(opts)

		// Check the error returned by the main function
		require.NoError(t, err)
	}()
	cancel()

	// // Wait for the main function to finish
	// select {
	// case <-time.After(5 * time.Second):
	// 	t.Fatal("Timeout: main function did not finish")
	// case <-ctx.Done():
	// }

	// // Assert that LoadConfig was called with the correct arguments
	// mockConfigLoader.AssertExpectations(t)
}

func mainWithOptions(opts Options) error {
	// Function similar to main but taking Options as an argument
	// This allows us to run the main function with different configurations in tests
	var configLoader config.ConfigLoader
	configLoader = &config.ConfigLoaderImpl{}
	cnf, err := configLoader.LoadConfig(opts.Config)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.Server{
		Listen:         opts.Listen,
		PinSize:        opts.PinSize,
		MaxExpire:      opts.MaxExpire,
		MaxPinAttempts: opts.MaxPinAttempts,
		WebRoot:        opts.WebRoot,
		Secret:         opts.Secret,
		Version:        revision,
		Config:         cnf,
	}

	return srv.Run(ctx)
}
