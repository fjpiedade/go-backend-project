package store

import (
	"context"
	"database/sql"
)

type PostStorageInterface interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
}

type UserStorageInterface interface {
	Create(context.Context, *User) error
	GetByID(context.Context, int64) (*User, error)
	GetAll(context.Context) ([]User, error)
	Update(context.Context, *User) error
	Delete(context.Context, int64) error
}

type Storage struct {
	Posts PostStorageInterface
	Users UserStorageInterface
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
