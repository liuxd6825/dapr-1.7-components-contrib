package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type AggregateRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Aggregate) error
	Delete(ctx context.Context, tenantId string, id string) error
	Update(ctx context.Context, tenantId string, v *model.Aggregate) error
	UpdateIsDelete(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error)
	FindById(ctx context.Context, tenantId string, id string) (*model.Aggregate, bool, error)
	DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error)
	NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, count uint64) (*model.Aggregate, bool, error)
}
