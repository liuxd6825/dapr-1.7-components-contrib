package service

import (
	"context"
	"errors"
	es "github.com/dapr/components-contrib/eventsourcing/v1"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateService interface {
	ExistAggregate(ctx context.Context, tenantId, aggregateId string) (*es.ExistAggregateResponse, error)
}

type aggregateService struct {
	repos *repository.AggregateRepository
}

func NewAggregateService(client *mongo.Client, collection *mongo.Collection) AggregateService {
	return &aggregateService{repos: repository.NewAggregateRepository(client, collection)}
}

func (c *aggregateService) ExistAggregate(ctx context.Context, tenantId, aggregateId string) (*es.ExistAggregateResponse, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId 不能为空")
	}
	if aggregateId == "" {
		return nil, errors.New("aggregateId 不能为空")
	}
	ok, err := c.repos.ExistAggregate(ctx, tenantId, aggregateId)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, err
	}
	return &es.ExistAggregateResponse{
		IsExist: ok,
	}, nil
}
