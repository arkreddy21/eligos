package eligos

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type UserServiceI interface {
	CreateUser(u *User) error
	GetUser(email string) (*User, error)
	GetUserById(id uuid.UUID) (*User, error)
	GetSpaces(userid uuid.UUID) (*[]Space, error)
}

type Space struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type SpaceServiceI interface {
	CreateSpace(space *Space, userid uuid.UUID) error
	AddUserById(userid, spaceid uuid.UUID) error
	RemoveUserById(userid, spaceid uuid.UUID) error
	GetUsersInSpace(spaceid uuid.UUID) (*[]User, error)
	DeleteSpaceById(spaceid uuid.UUID) error
}

type Message struct {
	Id        uuid.UUID          `json:"id"`
	UserId    uuid.UUID          `json:"userid"`
	SpaceId   uuid.UUID          `json:"spaceid"`
	Body      string             `json:"body"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type MessageServiceI interface {
	CreateMessage(m *Message) error
}
