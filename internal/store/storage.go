package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	User interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
		GetAll(context.Context) ([]*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Delete(context.Context, int64) error
	}
	Document interface {
		GetDocumentByDocID(context.Context, string) (*Document, error)
		CreateDocument(context.Context, string, string) error
		UpdateDocument(context.Context, string, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		User:     &UserStore{db},
		Document: &DocumentStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
