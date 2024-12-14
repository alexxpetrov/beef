package rooms

import (
	"app-websocket/internal/domain"
	"context"

	"github.com/gocql/gocql"
)

type RoomStorage interface {
	GetAllRooms(ctx context.Context) ([]domain.Room, error)
	GetRoom(ctx context.Context, roomID gocql.UUID) (*domain.Room, error)
	CreateRoom(ctx context.Context, name string) (*domain.Room, error)
}

type RoomProvider struct {
	storage RoomStorage
}

func New(storage RoomStorage) *RoomProvider {
	return &RoomProvider{
		storage: storage,
	}
}

func (r *RoomProvider) GetAllRooms(ctx context.Context) ([]domain.Room, error) {
	return r.storage.GetAllRooms(ctx)
}

func (r *RoomProvider) GetRoom(ctx context.Context, roomID gocql.UUID) (*domain.Room, error) {
	return r.storage.GetRoom(ctx, roomID)
}

func (r *RoomProvider) CreateRoom(ctx context.Context, name string) (*domain.Room, error) {
	return r.storage.CreateRoom(ctx, name)
}
