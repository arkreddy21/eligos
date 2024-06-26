package postgres

import (
	"context"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type MessageService struct {
	db *DB
}

func NewMessageService(db *DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) CreateMessage(m eligos.MessageWUser) (eligos.MessageWUser, error) {
	m.CreatedAt = time.Now()
	m.Id = uuid.New()
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO messages VALUES ($1, $2, $3, $4, $5)", m.Id, m.UserId, m.SpaceId, m.Body, m.CreatedAt)
	if err != nil {
		return eligos.MessageWUser{}, err
	}
	return m, nil
}

func (s *MessageService) GetMessages(spaceid uuid.UUID) (*[]eligos.MessageWUser, error) {
	rows, err := s.db.dbpool.Query(context.Background(), "SELECT messages.*, users.name, users.email FROM messages JOIN users ON messages.userid = users.id WHERE spaceid = $1 ORDER BY createdat", spaceid)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	messages, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (eligos.MessageWUser, error) {
		var message eligos.MessageWUser
		err := row.Scan(&message.Id, &message.UserId, &message.SpaceId, &message.Body, &message.CreatedAt, &message.User.Name, &message.User.Email)
		message.User.Id = message.UserId
		return message, err
	})
	if err != nil {
		return nil, err
	}
	return &messages, nil
}
