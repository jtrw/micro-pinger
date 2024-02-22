package sender

import "log"

type Message struct {
	Text        string
	Webhook     string
	Datetime    string
	StatusCode  int
	Url         string
	ServiceName string
}

type Sender interface {
	Send() error
}

func NewSender(senderType string, message Message) Sender {
	switch senderType {
	case "telegram":
		return NewTelegram(message)
	case "slack":
		return NewSlack(message)
	default:
		log.Printf("[%s] Unsupported alert type: %s", message.ServiceName, senderType)
	}
	return nil
}
