package message_cache

import (
	"app-websocket/internal/config"
	"app-websocket/internal/domain"
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

type ChatCache interface {
	GetLastMessagesFromRoom(ctx context.Context, roomID string, count int) ([]domain.Message, error)
	GetRoomClients(ctx context.Context, roomID string) ([]domain.User, error)
	AddRoomClient(ctx context.Context, roomID string, user *domain.User) error
	DeleteClient(ctx context.Context, roomID string, user *domain.User) error
}

type ChatPersistentStorage interface {
	GetLastMessagesFromRoom(ctx context.Context, roomID gocql.UUID, count int) ([]domain.Message, error)
}

type ChatCacheProvider struct {
	cache             ChatCache
	persistentStorage ChatPersistentStorage
	countMessagesGet  int
}

func New(config *config.ChatConfig, cache ChatCache, persistentStorage ChatPersistentStorage) *ChatCacheProvider {
	return &ChatCacheProvider{
		cache:             cache,
		countMessagesGet:  config.CountMessagesGet,
		persistentStorage: persistentStorage,
	}
}

func (c *ChatCacheProvider) GetLastMessagesFromRoom(ctx context.Context, roomID gocql.UUID) ([]domain.Message, error) {
	messages, err := c.cache.GetLastMessagesFromRoom(ctx, roomID.String(), c.countMessagesGet)

	if err != nil {
		messages, err = c.persistentStorage.GetLastMessagesFromRoom(ctx, roomID, c.countMessagesGet)
		if err != nil {
			return nil, fmt.Errorf("services.message_cache.GetLastMessagesFromRoom: %w", err)
		}

		return messages, nil
	}

	return messages, nil
}

func (c *ChatCacheProvider) GetRoomClients(ctx context.Context, roomID string) ([]domain.User, error) {
	return c.cache.GetRoomClients(ctx, roomID)
}

func (c *ChatCacheProvider) AddRoomClient(ctx context.Context, roomID string, user *domain.User) error {
	return c.cache.AddRoomClient(ctx, roomID, user)
}

func (c *ChatCacheProvider) DeleteClient(ctx context.Context, roomID string, user *domain.User) error {
	return c.cache.DeleteClient(ctx, roomID, user)
}
