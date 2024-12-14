package message_online

import (
	erdtreev1 "app-websocket/gen/erdtree/v1"
	"app-websocket/gen/erdtree/v1/erdtreev1connect"
	"app-websocket/internal/domain"
	"app-websocket/internal/ports/ws"
	"context"
	"fmt"

	"connectrpc.com/connect"
)

type MessagePusher interface {
	Produce(msg *domain.Message) error
}

type MessageConsumer interface {
	Consume(ctx context.Context, handler domain.MessageHandler) error
}

type RoomClientsStorage interface {
	AddRoomClient(ctx context.Context, roomID string, user *domain.User) error
	DeleteClient(ctx context.Context, roomID string, user *domain.User) error
}

type MessageOnlineService struct {
	pusher          MessagePusher
	consumer        MessageConsumer
	roomClients     RoomClientsStorage
	hub             *ws.Hub
	userInfoService erdtreev1connect.ErdtreeStoreClient
}

func New(pusher MessagePusher, consumer MessageConsumer, roomClients RoomClientsStorage, hub *ws.Hub, erdtreeClient erdtreev1connect.ErdtreeStoreClient) *MessageOnlineService {
	return &MessageOnlineService{
		pusher:          pusher,
		consumer:        consumer,
		roomClients:     roomClients,
		hub:             hub,
		userInfoService: erdtreeClient,
	}
}

func (m *MessageOnlineService) PushMessage(ctx context.Context, msg *domain.Message) error {
	messageSentRequest := connect.NewRequest(&erdtreev1.SetRequest{
		Key:   msg.UserID.String() + "-lastMessage",
		Value: []byte(msg.Content),
	})

	_, err := m.userInfoService.Set(ctx, messageSentRequest)

	if err != nil {
		fmt.Println(fmt.Errorf("service.MessageOnlineService.Subscribe: Erdtree Set %w", err))
	}

	messageRoomRequest := connect.NewRequest(&erdtreev1.SetRequest{
		Key:   msg.UserID.String() + "-lastMessageRoom",
		Value: []byte(msg.RoomID.String()),
	})

	_, err = m.userInfoService.Set(ctx, messageRoomRequest)

	if err != nil {
		fmt.Println(fmt.Errorf("service.MessageOnlineService.Subscribe: Erdtree Set %w", err))
	}

	return m.pusher.Produce(msg)
}

func (m *MessageOnlineService) Consume(ctx context.Context, handler func(message domain.Message) error) error {
	return m.consumer.Consume(ctx, handler)
}

func (m *MessageOnlineService) Subscribe(ctx context.Context, client *ws.Client) error {
	m.hub.AddConnection(client)

	err := m.roomClients.AddRoomClient(ctx, client.RoomID.String(), client.User)
	if err != nil {
		return fmt.Errorf("service.MessageOnlineService.Subscribe: %w", err)
	}

	// newId := uuid.New()
	// msgID, _ := gocql.ParseUUID(newId.String())

	joinRoomRequest := connect.NewRequest(&erdtreev1.SetRequest{
		Key:   client.User.ID.String() + "-joinedRoom",
		Value: []byte(client.RoomID.String()),
	})

	_, err = m.userInfoService.Set(ctx, joinRoomRequest)

	if err != nil {
		return fmt.Errorf("service.MessageOnlineService.Subscribe: Erdtree Set %w", err)
	}

	// return m.pusher.Produce(&domain.Message{
	// 	ID:          msgID,
	// 	Content:     "joined the room",
	// 	RoomID:      client.RoomID,
	// 	UserID:      client.User.ID,
	// 	TimeCreated: time.Now(),
	// 	Nickname:    client.User.Nickname,
	// })
	return nil

}

func (m *MessageOnlineService) Unsubscribe(ctx context.Context, client *ws.Client) error {
	m.hub.DeleteConnection(client)

	err := m.roomClients.DeleteClient(ctx, client.RoomID.String(), client.User)
	if err != nil {
		return fmt.Errorf("service.MessageOnlineService.Unsubscribe: %w", err)
	}

	// uuid := uuid.New()
	// msgId, _ := gocql.ParseUUID(uuid.String())

	joinRoomRequest := connect.NewRequest(&erdtreev1.SetRequest{
		Key:   client.User.ID.String() + "-leftRoom",
		Value: []byte(client.RoomID.String()),
	})

	_, err = m.userInfoService.Set(ctx, joinRoomRequest)

	if err != nil {
		return fmt.Errorf("service.MessageOnlineService.Subscribe: Erdtree Set %w", err)
	}

	// return m.pusher.Produce(&domain.Message{
	// 	ID:          msgId,
	// 	Content:     "left the room",
	// 	RoomID:      client.RoomID,
	// 	UserID:      client.User.ID,
	// 	TimeCreated: time.Now(),
	// 	Nickname:    client.User.Nickname,
	// })
	return nil
}
