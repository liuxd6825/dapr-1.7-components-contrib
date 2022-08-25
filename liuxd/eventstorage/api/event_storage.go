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
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
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

func NewEventStorage(log logger.Logger) eventstorage.EventStorage {
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
func (s *EventStorage) LoadEvent(ctx context.Context, req *dto.LoadEventRequest) (*dto.LoadResponse, error) {
	res, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		sequenceNumber := uint64(0)
		//获取最后的聚合镜像
		snapshot, ok, err := s.snapshotService.FindByMaxSequenceNumber(ctx, req.TenantId, req.AggregateId, req.AggregateType)
		if err != nil {
			return nil, newError("findByMaxSequenceNumber() error taking snapshot.", err)
		}
		if ok {
			sequenceNumber = snapshot.SequenceNumber
		}
		//获取聚合镜像之后的事件列表
		events, ok, err := s.eventService.FindBySequenceNumber(ctx, req.TenantId, req.AggregateId, req.AggregateType, sequenceNumber)
		if err != nil {
			return nil, newError("findBySequenceNumber() error taking events.", err)
		}
		return NewLoadResponse(req.TenantId, req.AggregateId, req.AggregateType, snapshot, events), err
	})
	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	if res != nil {
		resp, _ := res.(*dto.LoadResponse)
		resp.Headers = headers
		return resp, err
	}
	return &dto.LoadResponse{Headers: headers}, err
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
func (s *EventStorage) CreateEvent(ctx context.Context, req *dto.CreateEventRequest) (*dto.CreateEventResponse, error) {
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
	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = dto.ResponseStatusEventDuplicate
	}
	return &dto.CreateEventResponse{Headers: headers}, err
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
func (s *EventStorage) DeleteEvent(ctx context.Context, req *dto.DeleteEventRequest) (*dto.DeleteEventResponse, error) {
	isDuplicateEvent := false
	_, err := s.doSession(ctx, func(ctx context.Context) (interface{}, error) {
		agg, ok, err := s.aggregateService.SetDeleted(ctx, req.TenantId, req.AggregateId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is not fond ", req.AggregateId))
		}
		if agg.Deleted {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is deleted", req.AggregateId))
		}

		events := []*dto.EventDto{
			req.Event,
		}

		err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, events, agg.SequenceNumber+1)
		if err, isDuplicateEvent = NotDuplicateKeyError(err); err != nil {
			return nil, err
		}
		return nil, nil
	})
	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = dto.ResponseStatusEventDuplicate
	}
	return &dto.DeleteEventResponse{Headers: headers}, err
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
func (s *EventStorage) ApplyEvent(ctx context.Context, req *dto.ApplyEventsRequest) (*dto.ApplyEventsResponse, error) {
	isDuplicateEvent := false
	res, err := s.doSession(ctx, func(ctx context.Context) (any, error) {
		if req == nil {
			return nil, errors.New("request is nil")
		}
		length := len(req.Events)
		if length == 0 {
			return nil, errors.New("request.events size 0 ")
		}
		agg, ok, sn, err := s.aggregateService.NextSequenceNumber(ctx, req.TenantId, req.AggregateId, uint64(length))
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" not found", req.AggregateId))
		}
		if agg.Deleted {
			return nil, errors.New(fmt.Sprintf("aggregate id \"%s\" is already deleted.", req.AggregateId))
		}

		err = s.saveEvents(ctx, req.TenantId, req.AggregateId, req.AggregateType, req.Events, sn)
		/*		if err, isDuplicateEvent = NotDuplicateKeyError(err); err != nil {
				return nil, err
			}*/
		return nil, err
	})
	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	if isDuplicateEvent {
		headers.Status = dto.ResponseStatusEventDuplicate
	}
	if res != nil {
		resp := res.(*dto.ApplyEventsResponse)
		resp.Headers = headers
		return resp, err
	}
	return &dto.ApplyEventsResponse{Headers: headers}, err
}

func (s *EventStorage) FindRelations(ctx context.Context, req *dto.FindRelationsRequest) (*dto.FindRelationsResponse, error) {
	res, err := s.doSession(ctx, func(ctx context.Context) (any, error) {
		findRes, _, err := s.relationService.FindPaging(ctx, req)
		if err != nil {
			return nil, err
		}
		var errMsg string
		if findRes.Error != nil {
			errMsg = findRes.Error.Error()
		}
		var relations []*dto.Relation
		if findRes.Data != nil {
			for _, item := range findRes.Data {
				rel := dto.Relation{
					Id:            item.Id,
					TenantId:      item.TenantId,
					TableName:     item.TableName,
					AggregateId:   item.AggregateId,
					AggregateType: item.AggregateType,
					IsDeleted:     item.IsDeleted,
					RelValue:      item.RelValue,
					RelName:       item.RelName,
				}
				relations = append(relations, &rel)
			}
		}
		res := &dto.FindRelationsResponse{
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
	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	if res != nil {
		resp, _ := res.(*dto.FindRelationsResponse)
		resp.Headers = headers
		return resp, err
	}
	return &dto.FindRelationsResponse{Headers: headers}, err
}

func (s *EventStorage) FindEvents(ctx context.Context, req *dto.FindEventsRequest) (*dto.FindEventsResponse, error) {
	res := s.eventService.FindPaging(ctx, req)
	var errMessage string
	if res.Error != nil {
		errMessage = res.Error.Error()
	}
	resp := &dto.FindEventsResponse{
		Data:       NewFindEventsItems(res.Data),
		Headers:    dto.NewResponseHeadersSuccess(nil),
		TotalPages: res.TotalPages,
		TotalRows:  res.TotalRows,
		Filter:     res.Filter,
		Sort:       res.Sort,
		Error:      errMessage,
		IsFound:    res.IsFound,
	}
	return resp, res.Error
}

//
// RepublishMessage
// @Description:  发送消息列表中的事件
// @receiver s
// @param ctx
// @param req
// @return *eventstorage.RepublishMessageResponse
// @return error
//
func (s *EventStorage) RepublishMessage(ctx context.Context, req *dto.RepublishMessageRequest) (*dto.RepublishMessageResponse, error) {
	resp := &dto.RepublishMessageResponse{}
	limit := int64(100)
	if req.Limit > 0 {
		limit = req.Limit
	}
	list, ok, err := s.messageService.FindAll(ctx, &limit)
	if err != nil {
		return nil, err
	}
	if !ok {
		return resp, nil
	}
	for _, item := range list {
		if err := s.publishMessage(ctx, item.Event, true); err != nil {
			return resp, nil
		}
	}
	return resp, nil
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
func (s *EventStorage) SaveSnapshot(ctx context.Context, req *dto.SaveSnapshotRequest) (*dto.SaveSnapshotResponse, error) {
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

	headers := dto.NewResponseHeaders(dto.ResponseStatusSuccess, err, nil)
	return &dto.SaveSnapshotResponse{Headers: headers}, nil
}

func (s *EventStorage) DeleteAggregate(ctx context.Context, req *dto.DeleteAggregateRequest) (*dto.DeleteAggregateResponse, error) {
	_, err := s.doSession(ctx, func(ctx context.Context) (any, error) {
		if err := s.aggregateService.DeleteById(ctx, req.TenantId, req.AggregateId); err != nil {
			return nil, err
		}
		if err := s.eventService.DeleteByAggregateId(ctx, req.TenantId, req.AggregateId); err != nil {
			return nil, err
		}
		if err := s.relationService.DeleteByAggregateId(ctx, req.TenantId, req.AggregateId); err != nil {
			return nil, err
		}
		if err := s.messageService.DeleteByAggregateId(ctx, req.TenantId, req.AggregateId); err != nil {
			return nil, err
		}
		if err := s.snapshotService.DeleteByAggregateId(ctx, req.TenantId, req.AggregateId); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return &dto.DeleteAggregateResponse{}, err
}

func (s *EventStorage) saveEvents(ctx context.Context, tenantId string, aggregateId string, aggregateType string, eventDtoList []*dto.EventDto, startSequenceNumber uint64) error {
	if eventDtoList == nil {
		return errors.New("eventDtoList is nil")
	}
	length := len(eventDtoList)
	if length == 0 {
		return errors.New("request.saveEvents size 0 ")
	}

	for i := uint64(0); i < uint64(length); i++ {
		eventDto := eventDtoList[i]

		event := NewEvent(tenantId, aggregateId, aggregateType, startSequenceNumber+i, eventDto)
		relations := model.NewRelations(tenantId, event.EventId, event.EventType, aggregateId, aggregateType, eventDto.Relations)

		err := s.saveEvent(ctx, event, relations)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EventStorage) saveEvent(ctx context.Context, event *model.Event, relations []*model.Relation) error {
	// 创建新事件，并设置PublishStatus为Wait
	if _, err := s.createEvent(ctx, event); err != nil {
		return newError("createEvent() error saving event.", err)
	}

	// 通过领域事件，保存聚合关系
	if err := s.saveRelations(ctx, event.TenantId, relations); err != nil {
		return newError("relationService.Create() error.", err)
	}

	// 发送事件到消息队列，并设置 PublishStatus 为 PublishStatusSuccess
	if err := s.publishMessage(ctx, event, false); err != nil {
		return err
	}

	return nil
}

func (s *EventStorage) saveRelations(ctx context.Context, tenantId string, relation []*model.Relation) error {
	if err := s.relationService.CreateMany(ctx, tenantId, relation); err != nil {
		return newError("relationService.Create() error: ", err)
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
func (s *EventStorage) publishMessage(ctx context.Context, event *model.Event, isRepublish bool) error {
	if event == nil {
		return errors.New("publishMessage(ctx, event, isRepublish) error: event is nil")
	}
	tenantId := event.TenantId
	messageId := event.EventId

	// 不是重发
	if !isRepublish {
		message := &model.Message{
			Id:          messageId,
			AggregateId: event.AggregateId,
			TenantId:    tenantId,
			EventId:     event.EventId,
			CreateTime:  time.Now(),
			Event:       event,
		}
		if err := s.messageService.Create(ctx, message); err != nil {
			return err
		}
	}

	contentType := "json"
	publishData := dto.PublishData{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventVersion:   event.EventVersion,
		EventType:      event.EventType,
		SequenceNumber: event.SequenceNumber,
	}
	bytes, err := json.Marshal(publishData)
	if err != nil {
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
		return newError("publishMessage(ctx, event, isRepublish) error: failed to publish event.", err)
	}

	if err := s.messageService.Delete(ctx, tenantId, messageId); err != nil {
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
//  @return *Event
//  @return error
//
func (s *EventStorage) createEvent(ctx context.Context, event *model.Event) (*model.Event, error) {
	err := s.eventService.Create(ctx, event)
	return event, err
}

func (s *EventStorage) findEventById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return s.eventService.FindById(ctx, tenantId, id)
}

func newError(msgType string, err error) error {
	return errors.New(msgType + err.Error())
}

func (s *EventStorage) newAggregateEntity(req *dto.CreateEventRequest) (*model.Aggregate, error) {
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
