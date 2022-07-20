package es_mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/service"
	"github.com/liuxd6825/components-contrib/pubsub"
)

type EventStorage struct {
	mongodb          *db.MongoDB
	log              logger.Logger
	metadata         common.Metadata
	getPubsubAdapter eventstorage.GetPubsubAdapter
	eventService     service.EventService
	snapshotService  service.SnapshotService
	aggregateService service.AggregateService
}

// NewMongoEventSourcing 创建
func NewMongoEventSourcing(log logger.Logger) eventstorage.EventStorage {
	return &EventStorage{log: log, mongodb: db.NewMongoDB(log)}
}

//
// Init
// @Description: 初始化
// @receiver s
// @param metadata
// @param adapter
// @return error
//
func (s *EventStorage) Init(metadata common.Metadata, adapter eventstorage.GetPubsubAdapter) error {
	s.getPubsubAdapter = adapter
	s.metadata = metadata
	if err := s.mongodb.Init(metadata); err != nil {
		return err
	}

	storageMetadata := s.mongodb.StorageMetadata()
	aggregateCollection := s.mongodb.NewCollection(storageMetadata.AggregateCollectionName())
	eventCollection := s.mongodb.NewCollection(storageMetadata.EventCollectionName())
	snapshotCollection := s.mongodb.NewCollection(storageMetadata.SnapshotCollectionName())

	mongoClient := s.mongodb.GetClient()

	s.aggregateService = service.NewAggregateService(mongoClient, aggregateCollection)
	s.eventService = service.NewEventService(mongoClient, eventCollection)
	s.snapshotService = service.NewSnapshotService(mongoClient, snapshotCollection)

	return nil
}

//
// LoadEvent
// @Description: 加载聚合事件
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.LoadResponse
// @return error
//
func (s *EventStorage) LoadEvent(ctx context.Context, req *eventstorage.LoadEventRequest) (*eventstorage.LoadResponse, error) {
	sequenceNumber := uint64(0)
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

//
// CreateEvent
// @Description: 创建聚合根并应用领域事件
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.CreateEventResponse
// @return error
//
func (s *EventStorage) CreateEvent(ctx context.Context, req *eventstorage.CreateEventRequest) (*eventstorage.CreateEventResponse, error) {
	agg, err := s.aggregateService.FindById(ctx, req.TenantId, req.AggregateId)
	if err != nil {
		return nil, err
	}
	if agg != nil {
		return nil, errors.New(fmt.Sprintf("aggregateId \"%s\" already exists", req.AggregateId))
	}

	agg, err = s.newAggregateEntity(req)
	if err != nil {
		return nil, err
	}
	if err = s.aggregateService.Create(ctx, agg); err != nil {
		return nil, err
	}

	err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, req.Events, agg.SequenceNumber)
	if err != nil {
		return nil, err
	}

	return &eventstorage.CreateEventResponse{}, nil
}

//
// DeleteEvent
// @Description:
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.DeleteEventResponse
// @return error
//
func (s *EventStorage) DeleteEvent(ctx context.Context, req *eventstorage.DeleteEventRequest) (*eventstorage.DeleteEventResponse, error) {
	/*	agg, err := s.aggregateService.FindById(ctx, req.TenantId, req.AggregateId)
		if err != nil {
			return nil, err
		}
		if agg == nil {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" not found", req.AggregateId))
		}
		if agg.Deleted {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is deleted", req.AggregateId))
		}*/
	agg, err := s.aggregateService.DeleteAndNextSequenceNumber(ctx, req.TenantId, req.AggregateId)
	if err != nil {
		return nil, err
	}
	events := []eventstorage.EventDto{*req.Event}
	if err := s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, &events, agg.SequenceNumber+1); err != nil {
		return nil, err
	}
	return &eventstorage.DeleteEventResponse{}, nil
}

//
// ApplyEvent
// @Description: 应用多个领域事件
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.ApplyEventsResponse
// @return error
//
func (s *EventStorage) ApplyEvent(ctx context.Context, req *eventstorage.ApplyEventsRequest) (*eventstorage.ApplyEventsResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	length := len(*req.Events)
	if length == 0 {
		return nil, errors.New("request.events size 0 ")
	}
	agg, err := s.aggregateService.NextSequenceNumber(ctx, req.TenantId, req.AggregateId, uint64(length))
	if err != nil {
		return nil, err
	}
	if agg == nil {
		return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" not found", req.AggregateId))
	}
	if agg.Deleted {
		return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is already deleted.", req.AggregateId))
	}

	err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, req.Events, agg.SequenceNumber+1)
	if err != nil {
		return nil, err
	}
	return &eventstorage.ApplyEventsResponse{}, nil
}

func (s *EventStorage) saveEvents(ctx context.Context, tenantId string, aggregateId string, aggregateType string, events *[]eventstorage.EventDto, startSn uint64) error {
	if events == nil {
		return errors.New("events is nil")
	}
	length := len(*events)
	if length == 0 {
		return errors.New("request.saveEvents size 0 ")
	}

	var applyEvents []*eventstorage.Event
	list := *events
	for i := 0; i < length; i++ {
		appReq, err := eventstorage.NewEvent(tenantId, aggregateId, aggregateType, list[i])
		if err != nil {
			return err
		}
		applyEvents = append(applyEvents, appReq)
	}

	count := uint64(length)
	for i := uint64(0); i < count; i++ {
		applyEvent := applyEvents[i]
		err := s.saveEvent(ctx, applyEvent, startSn+i)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// SaveSnapshot
// @Description: 保存聚合镜像对象
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.SaveSnapshotResponse
// @return error
//
func (s *EventStorage) SaveSnapshot(ctx context.Context, req *eventstorage.SaveSnapshotRequest) (*eventstorage.SaveSnapshotResponse, error) {

	snapshot := &model.SnapshotEntity{
		Id:               model.NewObjectID(),
		TenantId:         req.TenantId,
		AggregateId:      req.AggregateId,
		AggregateType:    req.AggregateType,
		SequenceNumber:   req.SequenceNumber,
		Metadata:         req.Metadata,
		AggregateData:    req.AggregateData,
		AggregateVersion: req.AggregateVersion,
	}

	err := s.snapshotService.Create(ctx, snapshot)
	if err != nil {
		return nil, newError("SnapshotService.Create(). error saving snapshot.", err)
	}
	return &eventstorage.SaveSnapshotResponse{}, nil
}

func (s *EventStorage) saveEvent(ctx context.Context, req *eventstorage.Event, sequenceNumber uint64) error {
	// 创建新事件，并设置PublishStatus为Wait
	if _, err := s.createEvent(ctx, req, sequenceNumber); err != nil {
		return newError("createEvent() error saving event.", err)
	}

	// 发送事件到消息队列，并设置PublishStatus为PublishStatusSuccess
	if err := s.publishMessage(ctx, req); err != nil {
		return newError("publishMessage() failed to publish event.", err)
	}

	// 更新event的消息发送状态为成功
	if err := s.updatePublishStatue(ctx, req.EventId, eventstorage.PublishStatusSuccess); err != nil {
		return err
	}
	return nil
}

//
//  updatePublishStatue
//  @Description: 更新事件状态PublishStatus
//  @receiver s
//  @param ctx
//  @param event
//  @param publishStatue
//  @return error
//
func (s *EventStorage) updatePublishStatue(ctx context.Context, eventId string, publishStatue eventstorage.PublishStatus) error {
	return s.eventService.UpdatePublishStatue(ctx, eventId, publishStatue)
}

//
//  publishMessage
//  @Description: 发送事件到消息队列，并设置PublishStatus为PublishStatusSuccess
//  @receiver s
//  @param ctx 上下文
//  @param req
//  @param event
//  @return error
//
func (s *EventStorage) publishMessage(ctx context.Context, req *eventstorage.Event) error {
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

	if err := s.getPubsubAdapter().Publish(pubData); err != nil {
		return err
	}
	return nil
}

//
//  createEvent
//  @Description: 创建事件，并设置发送状态为PublishStatusWait
//  @receiver s
//  @param ctx
//  @param req
//  @return *EventEntity
//  @return error
//
func (s *EventStorage) createEvent(ctx context.Context, req *eventstorage.Event, sequenceNumber uint64) (*model.EventEntity, error) {
	idValue, err := model.ObjectIDFromHex(req.EventId)
	if err != nil {
		return nil, err
	}
	event := &model.EventEntity{
		Id:             idValue,
		TenantId:       req.TenantId,
		EventId:        req.EventId,
		Metadata:       req.Metadata,
		EventData:      req.EventData,
		EventVersion:   req.EventVersion,
		EventType:      req.EventType,
		AggregateId:    req.AggregateId,
		AggregateType:  req.AggregateType,
		PublishName:    req.PubsubName,
		Topic:          req.Topic,
		PublishStatus:  eventstorage.PublishStatusWait,
		SequenceNumber: sequenceNumber,
	}
	err = s.eventService.Create(ctx, event)
	return event, err
}

func (s *EventStorage) findEventById(ctx context.Context, tenantId string, id string) (*model.EventEntity, error) {
	return s.eventService.FindById(ctx, tenantId, id)
}

func newError(msgType string, err error) error {
	return errors.New(msgType + err.Error())
}

func (s *EventStorage) newAggregateEntity(req *eventstorage.CreateEventRequest) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(req.AggregateId)
	if err != nil {
		return nil, err
	}
	return &model.AggregateEntity{
		Id:             idValue,
		TenantId:       req.TenantId,
		AggregateId:    req.AggregateId,
		AggregateType:  req.AggregateType,
		SequenceNumber: 1,
	}, nil
}
