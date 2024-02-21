package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Telegram struct {
	Webhook string
}

type TelegramMessage struct {
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func NewTelegram(webhook string) Telegram {
	return Telegram{Webhook: webhook}
}

func (t Telegram) Send(message string) error {
	telegramMessage := TelegramMessage{Text: message, ParseMode: "html"}
	jsonMessage, err := json.Marshal(telegramMessage)
	if err != nil {
		return err
	}
	log.Printf("Telegram message: %v", string(jsonMessage))
	log.Printf("Telegram webhook: %v", t.Webhook)
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

	log.Printf("Telegram response: %v", resp)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return "ok", nil
}
