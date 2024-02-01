package postgres

import (
	"context"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SpaceService struct {
	db *DB
}

func NewSpaceService(db *DB) *SpaceService {
	return &SpaceService{db: db}
}

func (s *SpaceService) CreateSpace(space *eligos.Space, userid uuid.UUID) error {
	space.Id = uuid.New()
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO spaces VALUES ($1, $2)", space.Id, space.Name)
	if err != nil {
		return err
	}
	_, err = s.db.dbpool.Exec(context.Background(), "INSERT INTO userspaces VALUES ($1, $2)", userid, space.Id)
	return err
}

func (s *SpaceService) AddUserById(userid, spaceid uuid.UUID) error {
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO userspaces VALUES ($1, $2)", userid, spaceid)
	return err
}

func (s *SpaceService) GetUsersInSpace(spaceid uuid.UUID) (*[]eligos.User, error) {
	rows, err := s.db.dbpool.Query(context.Background(), "SELECT u.id, u.name, u.email FROM users u JOIN userspaces us ON u.id=us.userid WHERE us.spaceid=$1", spaceid)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (eligos.User, error) {
		var user eligos.User
		err := row.Scan(&user.Id, &user.Name, &user.Email)
		return user, err
	})
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *SpaceService) RemoveUserById(userid, spaceid uuid.UUID) error {
	_, err := s.db.dbpool.Exec(context.Background(), "DELETE FROM userspaces WHERE userid=$1 AND spaceid=$2", userid, spaceid)
	return err
}

func (s *SpaceService) DeleteSpaceById(spaceid uuid.UUID) error {
	_, err := s.db.dbpool.Exec(context.Background(), "DELETE FROM spaces WHERE id=$1", spaceid)
	return err
}
