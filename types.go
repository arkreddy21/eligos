package eligos

import (
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	Name     string
	Email    string
	Password string `json:"-"`
}

type UserServiceI interface {
	CreateUser(u *User) error
	GetUser(email string) (*User, error)
	GetUserById(id uuid.UUID) (*User, error)
}
