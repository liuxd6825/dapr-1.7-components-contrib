package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type SnapshotRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Snapshot) error
	DeleteById(ctx context.Context, tenantId string, id string) error
	DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error
	Update(ctx context.Context, tenantId string, v *model.Snapshot) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Snapshot, bool, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) ([]*model.Snapshot, bool, error)
	FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.Snapshot, bool, error)
}
