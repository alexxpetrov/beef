package pg

import (
	"app-websocket/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(pgURL string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("storage.pg.New: %w", err)
	}

	return &Postgres{
		pool: pool,
	}, nil
}

func (pg *Postgres) CloseConnection() {
	pg.pool.Close()
}

func (pg *Postgres) GetAllRooms(ctx context.Context) ([]domain.Room, error) {
	rows, err := pg.pool.Query(ctx, "SELECT id, name, time_created FROM rooms")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Room{}, nil
		}

		return nil, fmt.Errorf("storage.pg.GetRooms: %w", err)
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		var roomId uuid.UUID
		err = rows.Scan(&roomId, &room.Name, &room.TimeCreated)
		if err != nil {
			return nil, fmt.Errorf("storage.pg.GetRooms Scan Error: %w", err)
		}

		gocqlUUID, err := gocql.ParseUUID(roomId.String())

		if err != nil {
			return nil, fmt.Errorf("storage.pg.GetRoom: %w", err)
		}

		room.ID = gocqlUUID

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (pg *Postgres) CreateRoom(ctx context.Context, name string) (*domain.Room, error) {
	timeCreated := time.Now()
	newId := uuid.New()

	row := pg.pool.QueryRow(ctx, "INSERT INTO rooms(id, name, time_created) VALUES ($1, $2, $3) RETURNING id", newId, name, timeCreated)

	var id uuid.UUID
	err := row.Scan(&id)

	gocqlUUID := gocql.UUID(newId)

	if err != nil {
		return nil, fmt.Errorf("storage.pg.CreateRoom: %w", err)
	}

	return &domain.Room{
		ID:          gocqlUUID,
		Name:        name,
		TimeCreated: timeCreated,
	}, nil
}

func (pg *Postgres) GetRoom(ctx context.Context, roomID gocql.UUID) (*domain.Room, error) {
	pgRoomUUID, err := uuid.FromBytes(roomID.Bytes())

	if err != nil {
		return nil, fmt.Errorf("storage.pg.GetRoom: %w", err)
	}

	row := pg.pool.QueryRow(ctx, "SELECT name, time_created FROM rooms WHERE id = $1", pgRoomUUID)

	var room domain.Room
	err = row.Scan(&room.Name, &room.TimeCreated)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, fmt.Errorf("storage.pg.GetRoom: %w", err)
	}

	room.ID = roomID

	return &room, nil
}
