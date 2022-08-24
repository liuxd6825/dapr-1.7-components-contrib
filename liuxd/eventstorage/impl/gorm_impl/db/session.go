package db

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"gorm.io/gorm"
)

type session struct {
	db *gorm.DB
}

func NewSession(db *gorm.DB) eventstorage.Session {
	return &session{db: db}
}

func (r *session) UseTransaction(ctx context.Context, dbFunc eventstorage.SessionFunc) error {
	return r.db.Transaction(func(tx *gorm.DB) (err error) {
		defer func() {
			if e1 := recover(); e1 != nil {
				if e2, ok := e1.(error); ok {
					err = e2
				} else {
					err = fmt.Errorf("UseTransaction() error: %v ", e1)
				}
			}
		}()
		txCtx := NewContext(ctx, tx)
		return dbFunc(txCtx)
	})
}
