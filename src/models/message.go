package models

import "time"

type Message struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	User       *User     `json:"user"`
	ReceiverId int64     `json:"receiver_id"`
	Receiver   *User     `json:"receiver"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
