package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Slack struct {
	Message Message
}

type SlackMessage struct {
	Text string `json:"text"`
}

func NewSlack(message Message) Slack {
	return Slack{Message: message}
}

func (s Slack) Send() error {
	slackMessage := SlackMessage{Text: s.Message.Status}
	jsonMessage, err := json.Marshal(slackMessage)
	if err != nil {
		return err
	}

	_, err = s.post(s.Message.Webhook, jsonMessage)
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
