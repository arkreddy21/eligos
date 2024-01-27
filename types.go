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

type Space struct {
	Id   uuid.UUID
	Name string
}

type SpaceServiceI interface {
	CreateSpace(space *Space) error
	AddUserById(userid, spaceid uuid.UUID) error
	RemoveUserById(userid, spaceid uuid.UUID) error
	DeleteSpaceById(spaceid uuid.UUID) error
}
