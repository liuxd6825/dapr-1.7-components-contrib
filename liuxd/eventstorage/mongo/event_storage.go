package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/common"
	"github.com/dapr/components-contrib/liuxd/eventstorage"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/kit/logger"
)

type EventStorage struct {
	mongodb          *MongoDB
	log              logger.Logger
	metadata         common.Metadata
	getPubsubAdapter eventstorage.GetPubsubAdapter
	eventService     EventService
	snapshotService  SnapshotService
	aggregateService AggregateService
}

// NewMongoEventSourcing 创建
func NewMongoEventSourcing(log logger.Logger) eventstorage.EventStorage {
	return &EventStorage{log: log, mongodb: NewMongoDB(log)}
}

// Init 初始化
func (s *EventStorage) Init(metadata common.Metadata, adapter eventstorage.GetPubsubAdapter) error {
	s.getPubsubAdapter = adapter
	s.metadata = metadata
	if err := s.mongodb.Init(metadata); err != nil {
		return err
	}

	eventCollection := s.mongodb.NewCollection(s.mongodb.storageMetadata.eventCollectionName)
	snapshotCollection := s.mongodb.NewCollection(s.mongodb.storageMetadata.snapshotCollectionName)

	mongoClient := s.mongodb.GetClient()

	s.aggregateService = NewAggregateService(mongoClient, eventCollection)
	s.eventService = NewEventService(mongoClient, eventCollection)
	s.snapshotService = NewSnapshotService(mongoClient, snapshotCollection)

	return nil
}

func (s *EventStorage) ExistAggregate(ctx context.Context, req *eventstorage.ExistAggregateRequest) (*eventstorage.ExistAggregateResponse, error) {
	return s.aggregateService.ExistAggregate(ctx, req.TenantId, req.AggregateId)
}

// LoadEvents 加载聚合根对象
func (s *EventStorage) LoadEvents(ctx context.Context, req *eventstorage.LoadEventRequest) (*eventstorage.LoadResponse, error) {
	sequenceNumber := int64(0)
	snapshot, err := s.snapshotService.FindByMaxSequenceNumber(ctx, req.TenantId, req.AggregateId)
	if err != nil {
		return nil, newError("findByMaxSequenceNumber() error taking snapshot.", err)
	}
	if snapshot != nil {
		sequenceNumber = snapshot.SequenceNumber
	}
	events, err := s.eventService.FindBySequenceNumber(ctx, req.TenantId, req.AggregateId, sequenceNumber)
	if err != nil {
		return nil, newError("findBySequenceNumber() error taking events.", err)
	}
	resp := NewLoadResponse(req.TenantId, req.AggregateId, snapshot, events)
	return resp, nil
}

// ApplyEvent 应用领域事件
func (s *EventStorage) ApplyEvent(ctx context.Context, req *eventstorage.ApplyEventRequest) (*eventstorage.ApplyResponse, error) {
	event, err := s.findEventById(req.TenantId, req.EventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("findEventById(); %s", err.Error()))
	}
	if event != nil {
		return nil, errors.New(fmt.Sprintf("eventId=%s ; event ID already exists", req.EventId))
	}
	if err := s.publishMessage(req); err != nil {
		return nil, newError("publishMessage() failed to publish event.", err)
	}
	if err := s.saveDomainEvent(req); err != nil {
		return nil, newError("saveDomainEvent() error saving event.", err)
	}
	return &eventstorage.ApplyResponse{}, nil
}

func (s *EventStorage) SaveSnapshot(ctx context.Context, req *eventstorage.SaveSnapshotRequest) (*eventstorage.SaveSnapshotResponse, error) {
	snapshot := &SnapshotEntity{
		TenantId:          req.TenantId,
		AggregateId:       req.AggregateId,
		AggregateType:     req.AggregateType,
		SequenceNumber:    req.SequenceNumber,
		Metadata:          req.Metadata,
		AggregateData:     req.AggregateData,
		AggregateRevision: req.AggregateRevision,
	}

	err := s.snapshotService.Save(ctx, snapshot)
	if err != nil {
		return nil, newError("save(). error saving snapshot.", err)
	}
	return &eventstorage.SaveSnapshotResponse{}, nil
}

func (s *EventStorage) updatePublishStateOk(tenantId string, eventId string) error {
	return s.eventService.UpdatePublishStateOk(context.Background(), tenantId, eventId)
}

func (s *EventStorage) updatePublishStateError(tenantId string, eventId string, err error) error {
	return s.eventService.UpdatePublishStateError(context.Background(), tenantId, eventId, err)
}

func (s *EventStorage) publishMessage(req *eventstorage.ApplyEventRequest) error {
	contentType := "json"
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	pubData := &pubsub.PublishRequest{
		PubsubName:  req.PubsubName,
		Topic:       req.Topic,
		Metadata:    req.Metadata,
		ContentType: &contentType,
		Data:        bytes,
	}
	return s.getPubsubAdapter().Publish(pubData)
}

func (s *EventStorage) saveDomainEvent(req *eventstorage.ApplyEventRequest) error {
	event := &EventEntity{
		Id:            req.EventId,
		TenantId:      req.TenantId,
		EventId:       req.EventId,
		Metadata:      req.Metadata,
		EventData:     req.EventData,
		EventRevision: req.EventRevision,
		EventType:     req.EventType,
		AggregateId:   req.AggregateId,
		AggregateType: req.AggregateType,
		PublishName:   req.PubsubName,
		Topic:         req.Topic,
	}
	err := s.eventService.Save(context.Background(), event)
	return err
}

func (s *EventStorage) findEventById(tenantId string, id string) (*EventEntity, error) {
	return s.eventService.FindById(context.Background(), tenantId, id)
}

func newError(msgType string, err error) error {
	return errors.New(msgType + err.Error())
}
