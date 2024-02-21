package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Slack struct {
	Webhook string
}

type SlackMessage struct {
	Text string `json:"text"`
}

func NewSlack(webhook string) Slack {
	return Slack{Webhook: webhook}
}

func (s Slack) Send(message string) error {
	slackMessage := SlackMessage{Text: message}
	jsonMessage, err := json.Marshal(slackMessage)
	if err != nil {
		return err
	}

	_, err = s.post(s.Webhook, jsonMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s Slack) post(url string, data []byte) (string, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return "ok", nil
}
