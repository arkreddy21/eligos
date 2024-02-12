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

func (s *MessageService) CreateMessage(m *eligos.Message) error {
	m.CreatedAt = time.Now()
	m.Id = uuid.New()
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO messages VALUES ($1, $2, $3, $4, $5)", m.Id, m.UserId, m.SpaceId, m.Body, m.CreatedAt)
	return err
}

func (s *MessageService) GetMessages(spaceid uuid.UUID) (*[]eligos.Message, error) {
	rows, err := s.db.dbpool.Query(context.Background(), "SELECT * FROM messages WHERE spaceid = $1 ORDER BY createdat", spaceid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (eligos.Message, error) {
		var message eligos.Message
		err := row.Scan(&message.Id, &message.UserId, &message.SpaceId, &message.Body, &message.CreatedAt)
		return message, err
	})
	if err != nil {
		return nil, err
	}
	return &messages, nil
}
