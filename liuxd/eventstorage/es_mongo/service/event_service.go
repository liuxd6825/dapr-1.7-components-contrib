package service

import (
	"context"
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventService interface {
	Create(ctx context.Context, event *model.EventEntity) error
	Update(ctx context.Context, event *model.EventEntity) error
	FindById(ctx context.Context, tenantId string, id string) (*model.EventEntity, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*[]model.EventEntity, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) (*[]model.EventEntity, error)
	UpdatePublishStatue(ctx context.Context, eventId string, publishStatue eventstorage.PublishStatus) error
}

func NewEventService(mongodb *db.MongoDB, collection *mongo.Collection) EventService {
	return &eventService{repos: repository.NewEventRepository(mongodb, collection)}
}

type eventService struct {
	repos *repository.EventRepository
}

func (s *eventService) Update(ctx context.Context, event *model.EventEntity) error {
	if err := s.validation(event); err != nil {
		return err
	}
	return s.repos.Insert(ctx, event)
}

func (s *eventService) Create(ctx context.Context, event *model.EventEntity) error {
	if err := s.validation(event); err != nil {
		return err
	}
	if event.SequenceNumber < 0 {
		return errors.New("event.SequenceNumber is 0")
	}
	event.TimeStamp = utils.NewMongoNow()
	return s.repos.Insert(ctx, event)
}

func (s *eventService) FindById(ctx context.Context, tenantId string, id string) (*model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId 不能为空")
	}
	if id == "" {
		return nil, errors.New("aggregateId 不能为空")
	}
	return s.repos.FindById(ctx, tenantId, id)
}

func (s *eventService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*[]model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId 不能为空")
	}
	if aggregateId == "" {
		return nil, errors.New("aggregateId 不能为空")
	}
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId, aggregateType)
}

func (s *eventService) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) (*[]model.EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId 不能为空")
	}
	return s.repos.FindBySequenceNumber(ctx, tenantId, aggregateId, aggregateType, sequenceNumber)
}

func (s *eventService) UpdatePublishStatue(ctx context.Context, eventId string, publishStatue eventstorage.PublishStatus) error {
	return s.repos.UpdatePublishStatue(ctx, eventId, publishStatue)
}

func (s *eventService) validation(event *model.EventEntity) error {
	if event == nil {
		return errors.New("event is nil")
	}
	if event.Id == model.NilObjectID {
		return errors.New("event.id is empty")
	}
	if event.TenantId == "" {
		return errors.New("event.tenantId is empty")
	}
	if event.EventId == "" {
		return errors.New("event.eventId is empty")
	}
	if event.EventVersion == "" {
		return errors.New("event.eventRevision is empty")
	}
	if event.Topic == "" {
		return errors.New("event.topic is empty")
	}
	if event.AggregateType == "" {
		return errors.New("event.aggregateType is empty")
	}
	if event.AggregateId == "" {
		return errors.New("event.aggregateId is empty")
	}
	if event.PublishName == "" {
		return errors.New("event.publishName is empty")
	}
	return nil
}
