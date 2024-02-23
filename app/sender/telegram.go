package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Telegram struct {
	Message Message
}

type TelegramMessage struct {
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewTelegram(message Message) Telegram {
	return Telegram{Message: message}
}

func (t Telegram) Send() error {

	msg := getTextMessage(t.Message)

	telegramMessage := TelegramMessage{Text: msg, ParseMode: "markdown"}
	jsonMessage, err := json.Marshal(telegramMessage)
	if err != nil {
		return err
	}

	_, err = t.post(t.Message.Webhook, jsonMessage)
	if err != nil {
		return err
	}

	return nil
}

func (t Telegram) post(url string, data []byte) (string, error) {
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
