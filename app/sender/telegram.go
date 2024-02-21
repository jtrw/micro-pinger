package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Telegram struct {
	Webhook string
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewTelegram(webhook string) Telegram {
	return Telegram{Webhook: webhook}
}

func (t Telegram) Send(message string) error {
	telegramMessage := TelegramMessage{ChatID: "@micro_manager", Text: message, ParseMode: "Markdown"}
	jsonMessage, err := json.Marshal(telegramMessage)
	if err != nil {
		return err
	}

	_, err = t.post(t.Webhook, jsonMessage)
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
