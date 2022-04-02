package mongo

import (
	"context"
	"errors"
	"github.com/dapr/components-contrib/liuxd/common"
	es "github.com/dapr/components-contrib/liuxd/eventstorage"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateService interface {
	ExistAggregate(ctx context.Context, tenantId, aggregateId string) (*es.ExistAggregateResponse, error)
}

type aggregateService struct {
	repos *AggregateRepository
}

func NewAggregateService(client *mongo.Client, collection *mongo.Collection) AggregateService {
	return &aggregateService{repos: NewAggregateRepository(client, collection)}
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

type EventService interface {
	Save(ctx context.Context, event *EventEntity) error
	FindById(ctx context.Context, tenantId string, id string) (*EventEntity, error)
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]EventEntity, error)
	FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, sequenceNumber int64) (*[]EventEntity, error)
	UpdatePublishStateOk(ctx context.Context, tenantId string, id string) error
	UpdatePublishStateError(ctx context.Context, tenantId string, id string, err error) error
}

func NewEventService(client *mongo.Client, collection *mongo.Collection) EventService {
	return &eventService{repos: NewEventRepository(client, collection)}
}

type eventService struct {
	repos *EventRepository
}

func (s *eventService) Save(ctx context.Context, event *EventEntity) error {
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
	event.TimeStamp = common.NewMongoNow()
	event.SequenceNumber = s.repos.NextSequenceNumber(ctx, event.TenantId, event.AggregateId, event.AggregateType)
	return s.repos.Insert(ctx, event)
}

func (s *eventService) FindById(ctx context.Context, tenantId string, id string) (*EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId不能为空")
	}
	if id == "" {
		return nil, errors.New("aggregateId不能为空")
	}
	return s.repos.FindById(ctx, tenantId, id)
}

func (s *eventService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]EventEntity, error) {
	if tenantId == "" {
		return nil, errors.New("tenantId不能为空")
	}
	if aggregateId == "" {
		return nil, errors.New("aggregateId不能为空")
	}
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *eventService) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, sequenceNumber int64) (*[]EventEntity, error) {
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
	return s.repos.UpdatePublishState(ctx, tenantId, id, 1, common.NewMongoNow(), nil)
}

func (s *eventService) UpdatePublishStateError(ctx context.Context, tenantId string, id string, err error) error {
	if tenantId == "" {
		return errors.New("tenantId不能为空")
	}
	if id == "" {
		return errors.New("id不能为空")
	}
	return s.repos.UpdatePublishState(ctx, tenantId, id, -1, common.NewMongoNow(), err)
}

type SnapshotService interface {
	Save(ctx context.Context, snapshot *SnapshotEntity) error
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]SnapshotEntity, error)
	FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*SnapshotEntity, error)
}

func NewSnapshotService(client *mongo.Client, collection *mongo.Collection) SnapshotService {
	return &snapshotService{
		repos: NewSnapshotRepository(client, collection),
	}
}

type snapshotService struct {
	repos *SnapshotRepository
}

func (s *snapshotService) Save(ctx context.Context, snapshot *SnapshotEntity) error {
	snapshot.TimeStamp = common.NewMongoNow()
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]SnapshotEntity, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*SnapshotEntity, error) {
	return s.repos.FindByMaxSequenceNumber(ctx, tenantId, aggregateId)
}
