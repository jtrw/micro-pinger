package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	tempFile, err := ioutil.TempFile("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write sample YAML content to the temporary file
	yamlContent := []byte(`
services:
  - name: example
    url: https://example.com
    method: GET
    type: http
    body: ""
    interval: 5s
    headers:
      - name: "Content-Type"
        value: "application/json"
    response:
      status: 200
      body: "OK"
    alerts:
      - name: devops
        type: email
        webhook: "https://example.com/webhook"
        to: "devops@example.com"
        failure: 3
        success: 3
        send-on-resolve: true
`)
	_, err = tempFile.Write(yamlContent)
	assert.NoError(t, err)

	// Load configuration from the temporary file
	configFile := tempFile.Name()
	config, err := LoadConfig(configFile)

	// Assertions based on your expectations
	assert.NoError(t, err, "Expected no error loading config")
	assert.Len(t, config.Service, 1, "Expected one service in the config")

	// Add more assertions based on your specific YAML content and expectations
	// For example, check values of config.Service[0].Name, config.Service[0].URL, etc.
}
