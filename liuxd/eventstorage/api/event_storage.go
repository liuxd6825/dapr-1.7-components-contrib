package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/service"
	"github.com/liuxd6825/components-contrib/pubsub"
	"time"
)

type EventStorage struct {
	log              logger.Logger
	metadata         common.Metadata
	pubsubAdapter    eventstorage.GetPubsubAdapter
	eventService     service.EventService
	snapshotService  service.SnapshotService
	aggregateService service.AggregateService
	relationService  service.RelationService
	messageService   service.MessageService
	snapshotCount    uint64
	session          eventstorage.Session
}

func NewEventStorage(log logger.Logger) *EventStorage {
	return &EventStorage{
		log: log,
	}
}

func (s *EventStorage) Init(opts *eventstorage.Options) error {
	s.metadata = opts.Metadata
	s.pubsubAdapter = opts.PubsubAdapter
	s.eventService = service.NewEventService(opts.EventRepos.(repository.EventRepository))
	s.snapshotService = service.NewSnapshotService(opts.SnapshotRepos.(repository.SnapshotRepository))
	s.aggregateService = service.NewAggregateService(opts.AggregateRepos.(repository.AggregateRepository))
	s.relationService = service.NewRelationService(opts.RelationRepos.(repository.RelationRepository))
	s.messageService = service.NewMessageService(opts.MessageRepos.(repository.MessageRepository))
	s.snapshotCount = opts.SnapshotCount
	s.session = opts.Session
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
	res, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		sequenceNumber := uint64(0)
		snapshot, ok, err := s.snapshotService.FindByMaxSequenceNumber(ctx, req.TenantId, req.AggregateId, req.AggregateType)
		if err != nil {
			return nil, newError("findByMaxSequenceNumber() error taking snapshot.", err)
		}
		if ok {
			sequenceNumber = snapshot.SequenceNumber
		}
		events, ok, err := s.eventService.FindBySequenceNumber(ctx, req.TenantId, req.AggregateId, req.AggregateType, sequenceNumber)
		if err != nil {
			return nil, newError("findBySequenceNumber() error taking events.", err)
		}
		return NewLoadResponse(req.TenantId, req.AggregateId, req.AggregateType, snapshot, events), err
	})
	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	if res != nil {
		resp, _ := res.(*eventstorage.LoadResponse)
		resp.Headers = headers
		return resp, err
	}
	return &eventstorage.LoadResponse{Headers: headers}, err
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
	isDuplicateEvent := false
	_, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		agg, ok, err := s.aggregateService.FindById(ctx, req.TenantId, req.AggregateId)
		if err != nil {
			return nil, err
		}
		if ok {
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
		if err, isDuplicateEvent = NotDuplicateKeyError(err); err != nil {
			return nil, err
		}
		return nil, nil
	})
	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = eventstorage.ResponseStatusEventDuplicate
	}
	return &eventstorage.CreateEventResponse{Headers: headers}, err
}

//
// DeleteEvent
// @Description: 发布删除状态的事件
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.DeleteEventResponse
// @return error
//
func (s *EventStorage) DeleteEvent(ctx context.Context, req *eventstorage.DeleteEventRequest) (*eventstorage.DeleteEventResponse, error) {
	isDuplicateEvent := false
	_, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		agg, ok, err := s.aggregateService.SetIsDelete(ctx, req.TenantId, req.AggregateId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is not fond ", req.AggregateId))
		}
		if agg.Deleted {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is deleted", req.AggregateId))
		}

		events := []eventstorage.EventDto{*req.Event}

		err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, &events, agg.SequenceNumber+1)
		if err, isDuplicateEvent = NotDuplicateKeyError(err); err != nil {
			return nil, err
		}
		return nil, nil
	})
	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = eventstorage.ResponseStatusEventDuplicate
	}
	return &eventstorage.DeleteEventResponse{Headers: headers}, err
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
	isDuplicateEvent := false
	res, err := s.doSession(ctx, func(ctx context.Context) (any, error) {
		if req == nil {
			return nil, errors.New("request is nil")
		}
		length := len(*req.Events)
		if length == 0 {
			return nil, errors.New("request.events size 0 ")
		}
		agg, ok, err := s.aggregateService.NextSequenceNumber(ctx, req.TenantId, req.AggregateId, uint64(length))
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" not found", req.AggregateId))
		}
		if agg.Deleted {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is already deleted.", req.AggregateId))
		}

		err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, req.Events, agg.SequenceNumber)
		if err, isDuplicateEvent = NotDuplicateKeyError(err); err != nil {
			return nil, err
		}
		return nil, err
	})
	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = eventstorage.ResponseStatusEventDuplicate
	}
	if res != nil {
		resp := res.(*eventstorage.ApplyEventsResponse)
		resp.Headers = headers
		return resp, err
	}
	return &eventstorage.ApplyEventsResponse{Headers: headers}, err
}

func (s *EventStorage) GetRelations(ctx context.Context, req *eventstorage.GetRelationsRequest) (*eventstorage.GetRelationsResponse, error) {
	res, err := s.doSession(ctx, func(ctx context.Context) (any, error) {
		findRes, _, err := s.relationService.FindPaging(ctx, req)
		if err != nil {
			return nil, err
		}
		var errMsg string
		if findRes.Error != nil {
			errMsg = findRes.Error.Error()
		}
		var relations []*eventstorage.Relation
		if findRes.Data != nil {
			for _, item := range findRes.Data {
				rel := eventstorage.Relation{
					Id:          item.Id,
					TenantId:    item.TenantId,
					TableName:   item.TableName,
					AggregateId: item.AggregateId,
					IsDeleted:   item.IsDeleted,
					Items:       item.Items,
				}
				relations = append(relations, &rel)
			}
		}
		res := &eventstorage.GetRelationsResponse{
			Data:       relations,
			TotalRows:  uint64(findRes.TotalRows),
			TotalPages: uint64(findRes.TotalPages),
			PageSize:   uint64(findRes.PageSize),
			PageNum:    uint64(findRes.PageNum),
			Filter:     findRes.Filter,
			Sort:       findRes.Sort,
			Error:      errMsg,
		}
		return res, nil
	})
	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	if res != nil {
		resp, _ := res.(*eventstorage.GetRelationsResponse)
		resp.Headers = headers
		return resp, err
	}
	return &eventstorage.GetRelationsResponse{Headers: headers}, err
}

func (s *EventStorage) saveEvents(ctx context.Context, tenantId string, aggregateId string, aggregateType string, events *[]eventstorage.EventDto, startSequenceNumber uint64) error {
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
		err := s.saveEvent(ctx, applyEvent, startSequenceNumber+i)
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
	_, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		snapshot := &model.Snapshot{
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
		return nil, nil
	})

	headers := eventstorage.NewResponseHeaders(eventstorage.ResponseStatusSuccess, err, nil)
	return &eventstorage.SaveSnapshotResponse{Headers: headers}, nil
}

func (s *EventStorage) saveEvent(ctx context.Context, req *eventstorage.Event, sequenceNumber uint64) error {
	// 创建新事件，并设置PublishStatus为Wait
	if _, err := s.createEvent(ctx, req, sequenceNumber); err != nil {
		return newError("createEvent() error saving event.", err)
	}

	// 通过领域事件，保存聚合关系
	if err := s.saveRelations(ctx, req); err != nil {
		return newError("relationService.Create() error.", err)
	}

	// 发送事件到消息队列，并设置 PublishStatus 为 PublishStatusSuccess
	if err := s.publishMessage(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *EventStorage) saveRelations(ctx context.Context, req *eventstorage.Event) error {
	if req != nil && len(req.Relations) > 0 {
		relation := model.NewRelationEntity(req.TenantId, req.AggregateId, req.AggregateType, req.Relations)
		if err := s.relationService.Save(ctx, relation); err != nil {
			return newError("relationService.Create() error: ", err)
		}
	}
	return nil
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
func (s *EventStorage) publishMessage(ctx context.Context, event *eventstorage.Event) error {
	contentType := "json"
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &model.Message{
		Id:          event.EventId,
		AggregateId: event.AggregateId,
		TenantId:    event.TenantId,
		EventId:     event.EventId,
		CreateTime:  time.Now(),
		Event:       event,
	}

	if err := s.messageService.Create(ctx, message); err != nil {
		return err
	}

	pubData := &pubsub.PublishRequest{
		PubsubName:  event.PubsubName,
		Topic:       event.Topic,
		Metadata:    event.Metadata,
		ContentType: &contentType,
		Data:        bytes,
	}

	if err := s.pubsubAdapter().Publish(pubData); err != nil {
		return newError("publishMessage() error: failed to publish event.", err)
	}

	/*	if err := s.messageService.Delete(ctx, message.TenantId, message.Id); err != nil {
		return err
	}*/

	return nil
}

//
//  createEvent
//  @Description: 创建事件，并设置发送状态为PublishStatusWait
//  @receiver s
//  @param ctx
//  @param req
//  @return *Event
//  @return error
//
func (s *EventStorage) createEvent(ctx context.Context, req *eventstorage.Event, sequenceNumber uint64) (*model.Event, error) {
	event := &model.Event{
		Id:             req.EventId,
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
	err := s.eventService.Create(ctx, event)
	return event, err
}

func (s *EventStorage) findEventById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return s.eventService.FindById(ctx, tenantId, id)
}

func newError(msgType string, err error) error {
	return errors.New(msgType + err.Error())
}

func (s *EventStorage) newAggregateEntity(req *eventstorage.CreateEventRequest) (*model.Aggregate, error) {
	return &model.Aggregate{
		Id:             req.AggregateId,
		TenantId:       req.TenantId,
		AggregateId:    req.AggregateId,
		AggregateType:  req.AggregateType,
		SequenceNumber: 1,
	}, nil
}

func (s *EventStorage) doSession(ctx context.Context, fun func(ctx context.Context) (any, error)) (any, error) {
	var data interface{}
	var err error
	err = s.session.UseTransaction(ctx, func(ctx context.Context) error {
		data, err = fun(ctx)
		return err
	})
	return data, err
}

func (s *EventStorage) GetLogger() logger.Logger {
	return s.log
}
