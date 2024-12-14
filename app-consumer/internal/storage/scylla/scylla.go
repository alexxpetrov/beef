package scylla

import (
	"app-consumer/internal/domain"
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type Scylla struct {
	session *gocql.Session
}

func New(scyllaURL string) (*Scylla, error) {
	// Connect to ScyllaDB
	cluster := gocql.NewCluster(scyllaURL) // Replace with your container IP/hostname
	cluster.Keyspace = "messages_space"
	cluster.Consistency = gocql.LocalOne

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ScyllaDB: %v", err)
	}
	fmt.Println("storage.scylla.New Success")
	return &Scylla{
		session: session,
	}, nil
}

func (scylla *Scylla) Close() {
	if scylla.session != nil {
		scylla.session.Close()
		fmt.Println("ScyllaDB connection closed.")
	}
}

func (scylla *Scylla) PushMessage(ctx context.Context, msg *domain.Message) error {
	fmt.Println("storage.scylla.PushMessage", msg)
	id := uuid.New()
	gocqlUUID := gocql.UUID(id)

	if err := scylla.session.Query(
		`INSERT INTO messages (id, user_id, content, room_id, time_created) VALUES (?, ?, ?, ?, ?)`,
		gocqlUUID, msg.UserID, msg.Content, msg.RoomID, msg.TimeCreated,
	).WithContext(ctx).Exec(); err != nil {
		fmt.Println("ERROR?", err)
		return fmt.Errorf("storage.scylla.PushMessage: %w", err)
	}

	return nil
}
