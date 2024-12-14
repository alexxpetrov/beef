package redis

import (
	"app-websocket/internal/domain"

	"github.com/gocql/gocql"
)

func mapToUsers(m map[string]string) ([]domain.User, error) {
	var users []domain.User
	for key, value := range m {
		keyUuid, err := gocql.ParseUUID(key)

		if err != nil {
			return users, err
		}
		user := domain.User{
			ID:       keyUuid,
			Nickname: value,
		}

		users = append(users, user)
	}

	return users, nil
}
