package postgres

import (
	"context"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InviteService struct {
	db *DB
}

func NewInviteService(db *DB) *InviteService {
	return &InviteService{db: db}
}

func (s *InviteService) CreateInvite(invite *eligos.Invite) error {
	invite.Id = uuid.New()
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO invites VALUES ($1, $2, $3, $4)", invite.Id, invite.SpaceId, invite.SpaceName, invite.Email)
	return err
}

func (s *InviteService) DeleteInviteById(id uuid.UUID) error {
	_, err := s.db.dbpool.Exec(context.Background(), "DELETE FROM invites WHERE id=$1", id)
	return err
}

func (s *InviteService) GetInvitesByUser(email string) ([]eligos.Invite, error) {
	rows, err := s.db.dbpool.Query(context.Background(), "SELECT * FROM invites WHERE email=$1", email)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	invites, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (eligos.Invite, error) {
		var invite eligos.Invite
		err := row.Scan(&invite.Id, &invite.SpaceId, &invite.SpaceName, &invite.Email)
		return invite, err
	})
	if err != nil {
		return nil, err
	}
	return invites, nil
}
