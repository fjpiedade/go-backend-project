package store

import "errors"

var (
	ErrNotFound   = errors.New("record not found")
	ErrConflict   = errors.New("record already exists")
	ErrForeignKey = errors.New("related record not found")
)
