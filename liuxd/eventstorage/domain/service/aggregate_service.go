package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type AggregateService interface {
	Create(ctx context.Context, req *model.Aggregate) error
	DeleteById(ctx context.Context, tenantId, aggregateId string) error
	SetDeleted(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error)
	DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error)
	FindById(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error)
	NextSequenceNumber(ctx context.Context, tenantId, aggregateId string, count uint64) (*model.Aggregate, bool, uint64, error)
}

func NewAggregateService(repos repository.AggregateRepository) AggregateService {
	return &aggregateService{repos: repos}
}

type aggregateService struct {
	repos repository.AggregateRepository
}

func (c *aggregateService) DeleteById(ctx context.Context, tenantId, aggregateId string) error {
	return c.repos.DeleteByAggregateId(ctx, tenantId, aggregateId)
}

func (c *aggregateService) Destroy(ctx context.Context, tenantId, aggregateId string) error {
	return c.repos.DeleteById(ctx, tenantId, aggregateId)
}

func (c *aggregateService) Create(ctx context.Context, agg *model.Aggregate) error {
	if agg == nil {
		return nil
	}
	return c.repos.Create(ctx, agg)
}

func (c *aggregateService) SetDeleted(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	return c.repos.UpdateIsDelete(ctx, tenantId, aggregateId)
}

func (c *aggregateService) DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	return c.repos.DeleteAndNextSequenceNumber(ctx, tenantId, aggregateId)
}

func (c *aggregateService) FindById(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	return c.repos.FindById(ctx, tenantId, aggregateId)
}

func (c *aggregateService) NextSequenceNumber(ctx context.Context, tenantId, aggregateId string, count uint64) (*model.Aggregate, bool, uint64, error) {
	return c.repos.NextSequenceNumber(ctx, tenantId, aggregateId, count)
}
