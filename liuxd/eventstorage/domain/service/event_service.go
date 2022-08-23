package service

import (
	"context"
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
)

type EventService interface {
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, event *model.Event) error
	DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error)
	FindPaging(ctx context.Context, query dto.FindPagingQuery) *dto.FindPagingResult[*model.Event]
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, maxSequenceNumber uint64) ([]*model.Event, bool, error)
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

func (s *eventService) DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error {
	return s.repos.DeleteByAggregateId(ctx, tenantId, aggregateId)
}

func (s *eventService) FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return s.repos.FindById(ctx, tenantId, id)
}

func (s *eventService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId, aggregateType)
}

//
// FindBySequenceNumber 查找大于maxSequenceNumber的事件
func (s *eventService) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, maxSequenceNumber uint64) ([]*model.Event, bool, error) {
	return s.repos.FindByGtSequenceNumber(ctx, tenantId, aggregateId, aggregateType, maxSequenceNumber)
}

func (s *eventService) FindPaging(ctx context.Context, query dto.FindPagingQuery) *dto.FindPagingResult[*model.Event] {
	return s.repos.FindPaging(ctx, query)
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
	if event.PubsubName == "" {
		return errors.New("event.PubsubName is empty")
	}
	return nil
}
