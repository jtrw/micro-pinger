package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

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

	configFile := tempFile.Name()
	config, err := LoadConfig(configFile)

	assert.NoError(t, err, "Expected no error loading config")
	assert.Len(t, config.Service, 1, "Expected one service in the config")
}

func TestConfigNotFound(t *testing.T) {
	_, err := LoadConfig("non-existing-file.yaml")

	assert.Error(t, err, "Expected an error loading non-existing config file")
}

func TestBadConfig(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	yamlContent := []byte(`bad-yaml-content`)
	_, err = tempFile.Write(yamlContent)
	assert.NoError(t, err)

	configFile := tempFile.Name()
	_, err = LoadConfig(configFile)

	assert.Error(t, err, "Expected an error loading bad config file")
}
