package scylla

import (
	"app-websocket/internal/domain"
	"context"
	"fmt"

	"github.com/gocql/gocql"
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

func (scylla *Scylla) GetLastMessagesFromRoom(ctx context.Context, roomID gocql.UUID, count int) ([]domain.Message, error) {
	rows := scylla.session.Query(
		`SELECT m.content, u.nickname, m.user_id, m.time_created FROM messages AS m 
    		JOIN users AS u ON m.user_id = u.id 
            WHERE m.room_id = ?
            ORDER BY time_created
            LIMIT ?`, roomID, count).Iter()

	if rows.NumRows() == 0 {
		fmt.Println("No rows found")
		return []domain.Message{}, nil
	}
	defer rows.Close()

	var messages []domain.Message

	scanner := rows.Scanner()

	for scanner.Next() {
		msg := domain.Message{}

		scanner.Scan(msg)

		fmt.Println("storage.scylla.GetLastMessagesFromRoom", msg)

		messages = append(messages, msg)
	}
	fmt.Println("storage.scylla.GetLastMessagesFromRoom", messages)

	return messages, nil
}
