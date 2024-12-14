package chat

import (
	"time"

	"github.com/gocql/gocql"
)

type CreateRoomReq struct {
	Name string `json:"name"`
}

type RoomRes struct {
	ID          gocql.UUID `json:"id"`
	TimeCreated time.Time  `json:"time_created"`
	Name        string     `json:"name"`
}

type ClientRes struct {
	Username string `json:"nickname"`
}
