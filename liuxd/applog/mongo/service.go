package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventLogService interface {
	Insert(ctx context.Context, entity *EventLog) error
	Update(ctx context.Context, entity *EventLog) error
	FindById(ctx context.Context, tenantId string, id string) (*EventLog, error)
	FindBySubAppIdAndCommandId(ctx context.Context, tenantId string, subAppId string, commandId string) (*[]EventLog, error)
}

type eventLogService struct {
	repos *EventLogRepository
}

func NewEventLogService(client *mongo.Client, collection *mongo.Collection) EventLogService {
	return &eventLogService{
		repos: NewEventLogRepository(client, collection),
	}
}

func (e *eventLogService) Insert(ctx context.Context, entity *EventLog) error {
	return e.repos.Insert(ctx, entity)
}

func (e *eventLogService) Update(ctx context.Context, entity *EventLog) error {
	return e.repos.Update(ctx, entity)
}

func (e *eventLogService) FindBySubAppIdAndCommandId(ctx context.Context, tenantId string, subAppId string, commandId string) (*[]EventLog, error) {
	return e.repos.FindBySubAppIdAndCommandId(ctx, tenantId, subAppId, commandId)
}

func (e *eventLogService) FindById(ctx context.Context, tenantId string, id string) (*EventLog, error) {
	return e.repos.FindById(ctx, tenantId, id)
}

type AppLogService interface {
	Insert(ctx context.Context, entity *AppLog) error
	Update(ctx context.Context, entity *AppLog) error
	FindById(ctx context.Context, tenantId string, id string) (*AppLog, error)
}

type appLogService struct {
	repos *AppLogRepository
}

func NewAppLogService(client *mongo.Client, collection *mongo.Collection) AppLogService {
	return &appLogService{
		repos: NewAppLogRepository(client, collection),
	}
}

func (e *appLogService) Insert(ctx context.Context, entity *AppLog) error {
	return e.repos.Insert(ctx, entity)
}

func (e *appLogService) Update(ctx context.Context, entity *AppLog) error {
	return e.repos.Update(ctx, entity)
}

func (e *appLogService) FindById(ctx context.Context, tenantId string, id string) (*AppLog, error) {
	return e.repos.FindById(ctx, tenantId, id)
}
