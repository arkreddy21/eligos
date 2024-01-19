package postgres

import (
	"context"
	"github.com/arkreddy21/eligos"
	"github.com/google/uuid"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(u *eligos.User) error {
	u.Id = uuid.New()
	_, err := s.db.dbpool.Exec(context.Background(), "INSERT INTO users VALUES ($1, $2, $3, $4)", u.Id, u.Name, u.Email, u.Password)
	return err
}

func (s *UserService) GetUser(email string) (*eligos.User, error) {
	user := &eligos.User{}
	err := s.db.dbpool.QueryRow(context.Background(), "SELECT * FROM users WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserById(id uuid.UUID) (*eligos.User, error) {
	user := &eligos.User{}
	err := s.db.dbpool.QueryRow(context.Background(), "SELECT * FROM users WHERE id=$1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
