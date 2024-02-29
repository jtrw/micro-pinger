package sender

import (
	"fmt"
	"log"
)

type Message struct {
	Status      string
	Webhook     string
	Datetime    string
	Url         string
	ServiceName string
	Response    Response
}

type Response struct {
	Text string
	Err  error
	Code int
}

type Sender interface {
	Send() error
}

func NewSender(senderType string, message Message) (Sender, error) {
	switch senderType {
	case "telegram":
		return NewTelegram(message), nil
	case "slack":
		return NewSlack(message), nil
	default:
		log.Printf("[%s] Unsupported sender type: %s", message.ServiceName, senderType)
		return nil, fmt.Errorf("unsupported sender type: %s", senderType)
	}
	return nil, nil
}

func getTextMessage(message Message) string {
	if message.Response.Err != nil {
		return fmt.Sprintf("❗*Service:* %s\n*Status:* %s\n*Datetime:* %s\n*URL:* %s\n*Error:* %s",
			message.ServiceName, message.Status, message.Datetime, message.Url, message.Response.Err.Error())
	}

	return fmt.Sprintf("✅ *Service:* %s\n*Status:* %s\n*Datetime:* %s\n*URL:* %s",
		message.ServiceName, message.Status, message.Datetime, message.Url)
}
