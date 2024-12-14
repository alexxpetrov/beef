package ws

import (
	"time"

	"github.com/gocql/gocql"
)

type Message struct {
	ID          gocql.UUID `json:"id"`
	Content     string     `json:"content"`
	RoomID      gocql.UUID `json:"room_id"`
	Username    string     `json:"nickname"`
	UserID      gocql.UUID `json:"user_id"`
	TimeCreated time.Time  `json:"time_created"`
}
