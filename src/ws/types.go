package ws

import (
	"time"

	"github.com/playneta/go-sessions/src/models"
)

type MessageEvent struct {
	From     string    `json:"from"`
	To       string    `json:"to"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"date_time"`
}

type MessageJoin struct {
	User string `json:"user"`
}

type MessageError struct {
	Error string `json:"error"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewMessageEvent(message models.Message) Event {
	data := MessageEvent{
		From:     message.User.Email,
		Text:     message.Text,
		DateTime: message.CreatedAt,
	}

	if message.Receiver != nil {
		data.To = message.Receiver.Email
	}

	return Event{
		Type: "message",
		Data: data,
	}
}
