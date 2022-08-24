package db

import (
	ctx "context"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ContextKey struct {
}

var dbContextKey = &ContextKey{}

func NewContext(parentCtx ctx.Context, tx *gorm.DB) ctx.Context {
	return context.WithValue(parentCtx, dbContextKey, tx)
}

func GetTransaction(ctx ctx.Context) *gorm.DB {
	db, ok := ctx.Value(dbContextKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}
