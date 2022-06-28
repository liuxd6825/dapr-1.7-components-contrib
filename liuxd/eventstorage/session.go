package eventstorage

import "context"

type Session interface {
	UseTransaction(context.Context, SessionFunc) error
}

type SessionFunc func(ctx context.Context) error
