package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type MessageRepository interface {
	Create(ctx context.Context, v *model.Message) error
	DeleteById(ctx context.Context, tenantId string, id string) error
	DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error
	Update(ctx context.Context, v *model.Message) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Message, bool, error)
	FindAll(ctx context.Context, limit *int64) ([]*model.Message, bool, error)
}
