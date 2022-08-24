package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
)

type EventRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Event) error
	DeleteById(ctx context.Context, tenantId string, id string) error
	DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error
	Update(ctx context.Context, tenantId string, v *model.Event) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error)
	FindPaging(ctx context.Context, query dto.FindPagingQuery) *dto.FindPagingResult[*model.Event]
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error)
	FindByGtSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error)
}
