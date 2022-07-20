package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateService interface {
	Create(ctx context.Context, req *model.AggregateEntity) error
	Delete(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error)
	DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error)
	FindById(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error)
	NextSequenceNumber(ctx context.Context, tenantId, aggregateId string, count uint64) (*model.AggregateEntity, error)
}

type aggregateService struct {
	repos *repository.AggregateRepository
}

func NewAggregateService(client *mongo.Client, collection *mongo.Collection) AggregateService {
	return &aggregateService{repos: repository.NewAggregateRepository(client, collection)}
}

func (c *aggregateService) Create(ctx context.Context, req *model.AggregateEntity) error {
	return c.repos.Insert(ctx, req)
}

func (c *aggregateService) Delete(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error) {
	return c.repos.Delete(ctx, tenantId, aggregateId)
}

func (c *aggregateService) DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error) {
	return c.repos.DeleteAndNextSequenceNumber(ctx, tenantId, aggregateId)
}

func (c *aggregateService) FindById(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error) {
	return c.repos.FindById(ctx, tenantId, aggregateId)
}

func (c *aggregateService) NextSequenceNumber(ctx context.Context, tenantId, aggregateId string, count uint64) (*model.AggregateEntity, error) {
	return c.repos.NextSequenceNumber(ctx, tenantId, aggregateId, count)
}
