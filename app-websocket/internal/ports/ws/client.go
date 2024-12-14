package ws

import (
	"app-websocket/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ServiceChatPusher interface {
	PushMessage(ctx context.Context, msg *domain.Message) error
	Unsubscribe(ctx context.Context, client *Client) error
}

type Client struct {
	Conn    *websocket.Conn
	Message chan *Message
	Logger  *slog.Logger
	RoomID  gocql.UUID
	User    *domain.User
	Pusher  ServiceChatPusher
}

func (c *Client) WriteMessage() {
	defer c.Close()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		fmt.Println(message)
		err := c.Conn.WriteJSON(message)
		if err != nil {
			c.Logger.Error("can not send message to client",
				slog.String("Username", c.User.Nickname),
				slog.String("RoomID", c.RoomID.String()),
				slog.String("ClientID", c.User.ID.String()),
				slog.String("error", err.Error()))
		}
	}
}

func (c *Client) ReadMessage(ctx context.Context) {
	defer func() {
		err := c.Pusher.Unsubscribe(ctx, c)
		if err != nil {
			c.Logger.Error("failed to Unsubscribe from room:", slog.String("error", err.Error()))
		}

		c.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger.Error("can not write message to client",
					slog.String("username", c.User.Nickname),
					slog.String("RoomID", c.RoomID.String()),
					slog.String("ClientID", c.User.ID.String()),
					slog.String("error", err.Error()))
			}
			break
		}

		newUuid := uuid.New()
		msgId, _ := gocql.ParseUUID(newUuid.String())

		msg := &domain.Message{
			ID:          msgId,
			Content:     string(m),
			RoomID:      c.RoomID,
			Nickname:    c.User.Nickname,
			UserID:      c.User.ID,
			TimeCreated: time.Now(),
		}

		err = c.Pusher.PushMessage(ctx, msg)
		if err != nil {
			c.Logger.Error("failed to push message:", slog.String("error", err.Error()))
		}
	}
}

func (c *Client) Close() {
	err := c.Conn.Close()
	if err != nil {
		c.Logger.Error("failed to close WebSocket connection:", slog.String("error", err.Error()))
	}
}
