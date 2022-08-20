package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type EventRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Event) error
	Delete(ctx context.Context, tenantId string, id string) error
	Update(ctx context.Context, tenantId string, v *model.Event) error
	UpdatePublishStatue(ctx context.Context, tenantId string, eventId string, publishStatue eventstorage.PublishStatus) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error)
	FindNotPublishStatusSuccess(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error)
}
