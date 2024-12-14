package domain

import (
	"time"

	"github.com/gocql/gocql"
)

type Message struct {
	ID          gocql.UUID
	Content     string
	Nickname    string
	TimeCreated time.Time
	RoomID      gocql.UUID
	UserID      gocql.UUID
}

type User struct {
	ID           gocql.UUID
	Nickname     string
	PasswordHash string
}

type Room struct {
	ID          gocql.UUID
	Name        string
	TimeCreated time.Time
}

type MessageHandler func(msg Message) error
