package service

import (
	"context"
	"errors"
	"github.com/dapr/components-contrib/eventsourcing/utils"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventService interface {
	Save(ctx context.Context, event *model.EventEntity) error
	FindById(ctx context.Context, tenantId string, id string) (*model.EventEntity, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.EventEntity, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, sequenceNumber int64) (*[]model.EventEntity, error)
	UpdatePublishStateOk(ctx context.Context, tenantId string, id string) error
	UpdatePublishStateError(ctx context.Context, tenantId string, id string, err error) error
}

func NewEventService(client *mongo.Client, collection *mongo.Collection) EventService {
	return &eventService{repos: repository.NewEventRepository(client, collection)}
}

type eventService struct {
	repos *repository.EventRepository
}

func (s *eventService) Save(ctx context.Context, event *model.EventEntity) error {
	if event.Id == "" {
		return errors.New("id不能为空")
	}
	if event.TenantId == "" {
		return errors.New("tenantId不能为空")
	}
	if event.EventId == "" {
		return errors.New("eventId不能为空")
	}
	if event.EventRevision == "" {
		return errors.New("eventRevision不能为空")
	}
	if event.Topic == "" {
		return errors.New("topic不能为空")
	}
	if event.AggregateType == "" {
		return errors.New("aggregateType不能为空")
	}
	if event.AggregateId == "" {
		return errors.New("aggregateId不能为空")
	}
	if event.PublishName == "" {
		return errors.New("publishName不能为空")
	}
	event.TimeStamp = utils.NewMongoNow()
	event.SequenceNumber = s.repos.NextSequenceNumber(ctx, event.TenantId, event.AggregateId, event.AggregateType)
	return s.repos.Insert(ctx, event)
}

func (s *eventService) FindById(ctx context.Context, tenantId string, id string) (*model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId不能为空")
	}
	if id == "" {
		return nil, errors.New("aggregateId不能为空")
	}
	return s.repos.FindById(ctx, tenantId, id)
}

func (s *eventService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId不能为空")
	}
	if aggregateId == "" {
		return nil, errors.New("aggregateId不能为空")
	}
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *eventService) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, sequenceNumber int64) (*[]model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId不能为空")
	}
	return s.repos.FindBySequenceNumber(ctx, tenantId, aggregateId, sequenceNumber)
}

func (s *eventService) UpdatePublishStateOk(ctx context.Context, tenantId string, id string) error {
	if tenantId == "" {
		return errors.New("tenantId不能为空")
	}
	if id == "" {
		return errors.New("id不能为空")
	}
	return s.repos.UpdatePublishState(ctx, tenantId, id, 1, utils.NewMongoNow(), nil)
}

func (s *eventService) UpdatePublishStateError(ctx context.Context, tenantId string, id string, err error) error {
	if tenantId == "" {
		return errors.New("tenantId不能为空")
	}
	if id == "" {
		return errors.New("id不能为空")
	}
	return s.repos.UpdatePublishState(ctx, tenantId, id, -1, utils.NewMongoNow(), err)
}
