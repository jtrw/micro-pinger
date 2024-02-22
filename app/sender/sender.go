package sender

import "log"

type Message struct {
	Status      string
	Webhook     string
	Datetime    string
	Url         string
	ServiceName string
	ErrorMsg    ErrorMsg
}

type ErrorMsg struct {
	Text string
	Err  error
	Code int
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
