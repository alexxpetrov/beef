package domain

import "time"

type Message struct {
	ID          string
	Content     string
	Nickname    string
	TimeCreated time.Time
	RoomID      string
	UserID      string
}

type Room struct {
	ID          string
	Name        string
	TimeCreated time.Time
}

type MessageHandler func(msg Message) error
