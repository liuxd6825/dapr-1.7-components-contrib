package service

import (
	"context"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventLogService interface {
	Insert(ctx context.Context, entity *model.EventLog) error
	Update(ctx context.Context, entity *model.EventLog) error
	FindById(ctx context.Context, tenantId string, subAppId string, commandId string) (*model.EventLog, error)
}

type eventLogService struct {
	repos *repository.EventLogRepository
}

func NewEventLogService(client *mongo.Client, collection *mongo.Collection) EventLogService {
	return &eventLogService{
		repos: repository.NewEventLogRepository(client, collection),
	}
}

func (e *eventLogService) Insert(ctx context.Context, entity *model.EventLog) error {
	return e.repos.Insert(ctx, entity)
}

func (e *eventLogService) Update(ctx context.Context, entity *model.EventLog) error {
	return e.repos.Update(ctx, entity)
}

func (e *eventLogService) FindById(ctx context.Context, tenantId string, subAppId string, commandId string) (*model.EventLog, error) {
	return e.repos.FindById(ctx, tenantId, subAppId, commandId)
}
