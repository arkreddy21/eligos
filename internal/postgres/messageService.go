package postgres

import (
	"context"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
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
