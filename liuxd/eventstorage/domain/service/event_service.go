package service

import (
	"context"
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type EventService interface {
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, event *model.Event) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error)
	UpdatePublishStatue(ctx context.Context, tenantId string, eventId string, publishStatue eventstorage.PublishStatus) error
}

func NewEventService(repos repository.EventRepository) EventService {
	return &eventService{repos: repos}
}

type eventService struct {
	repos repository.EventRepository
}

func (s *eventService) Update(ctx context.Context, event *model.Event) error {
	if err := s.validation(event); err != nil {
		return err
	}
	return s.repos.Create(ctx, event.TenantId, event)
}

func (s *eventService) Create(ctx context.Context, event *model.Event) error {
	if err := s.validation(event); err != nil {
		return err
	}
	if event.SequenceNumber < 0 {
		return errors.New("event.SequenceNumber is 0")
	}
	event.TimeStamp = utils.NewMongoNow()
	return s.repos.Create(ctx, event.TenantId, event)
}

func (s *eventService) FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return s.repos.FindById(ctx, tenantId, id)
}

func (s *eventService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId, aggregateType)
}

func (s *eventService) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error) {
	return s.repos.FindBySequenceNumber(ctx, tenantId, aggregateId, aggregateType, sequenceNumber)
}

func (s *eventService) UpdatePublishStatue(ctx context.Context, tenantId string, eventId string, publishStatue eventstorage.PublishStatus) error {
	return s.repos.UpdatePublishStatue(ctx, tenantId, eventId, publishStatue)
}

func (s *eventService) validation(event *model.Event) error {
	if event == nil {
		return errors.New("event is nil")
	}
	if len(event.Id) == 0 {
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