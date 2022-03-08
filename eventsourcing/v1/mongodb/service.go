package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	es "github.com/dapr/components-contrib/eventsourcing/v1"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/kit/logger"
	"github.com/valyala/fasthttp"
)

type Service struct {
	mongodb          *MongoDB
	logger           logger.Logger
	metadata         es.Metadata
	getPubsubAdapter es.GetPubsubAdapter
}

// NewMongoEventSourcing 创建
func NewMongoEventSourcing(logger logger.Logger) es.EventSourcing {
	return &Service{logger: logger, mongodb: NewMongoDB(logger)}
}

// Init 初始化
func (c *Service) Init(metadata es.Metadata, adapter es.GetPubsubAdapter) error {
	c.getPubsubAdapter = adapter
	c.metadata = metadata
	if err := c.mongodb.Init(metadata); err != nil {
		return err
	}
	return nil
}

func (c *Service) ExistAggregate(ctx *fasthttp.RequestCtx, req *es.ExistAggregateRequest) (*es.ExistAggregateResponse, error) {
	return c.mongodb.aggregateService.ExistAggregate(ctx, req.TenantId, req.AggregateId)
}

// LoadEvents 加载聚合根对象
func (c *Service) LoadEvents(ctx *fasthttp.RequestCtx, req *es.LoadEventRequest) (*es.LoadResponse, error) {
	sequenceNumber := int64(0)
	snapshot, err := c.mongodb.snapshotService.FindByMaxSequenceNumber(ctx, req.TenantId, req.AggregateId)
	if err != nil {
		return nil, newError("findByMaxSequenceNumber() error taking snapshot.", err)
	}
	if snapshot != nil {
		sequenceNumber = snapshot.SequenceNumber
	}
	events, err := c.mongodb.eventService.FindBySequenceNumber(ctx, req.TenantId, req.AggregateId, sequenceNumber)
	if err != nil {
		return nil, newError("findBySequenceNumber() error taking events.", err)
	}
	resp := es.NewLoadResponse(req.TenantId, req.AggregateId, snapshot, events)
	return resp, nil
}

// ApplyEvent 应用领域事件
func (c *Service) ApplyEvent(ctx *fasthttp.RequestCtx, req *es.ApplyEventRequest) (*es.ApplyResponse, error) {
	event, err := c.findEventById(req.TenantId, req.EventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("findEventById(); %s", err.Error()))
	}
	if event != nil {
		return nil, errors.New(fmt.Sprintf("eventId=%s ; event ID already exists", req.EventId))
	}
	if err := c.publishMessage(req); err != nil {
		return nil, newError("publishMessage() failed to publish event.", err)
	}
	if err := c.saveDomainEvent(req); err != nil {
		return nil, newError("saveDomainEvent() error saving event.", err)
	}
	return &es.ApplyResponse{}, nil
}

func (c *Service) SaveSnapshot(ctx *fasthttp.RequestCtx, req *es.SaveSnapshotRequest) (*es.SaveSnapshotResponse, error) {
	snapshot := &model.SnapshotEntity{
		TenantId:          req.TenantId,
		AggregateId:       req.AggregateId,
		AggregateType:     req.AggregateType,
		SequenceNumber:    req.SequenceNumber,
		Metadata:          req.Metadata,
		AggregateData:     req.AggregateData,
		AggregateRevision: req.AggregateRevision,
	}

	err := c.mongodb.snapshotService.Save(ctx, snapshot)
	if err != nil {
		return nil, newError("save(). error saving snapshot.", err)
	}
	return &es.SaveSnapshotResponse{}, nil
}

func (c *Service) updatePublishStateOk(tenantId string, eventId string) error {
	return c.mongodb.eventService.UpdatePublishStateOk(context.Background(), tenantId, eventId)
}

func (c *Service) updatePublishStateError(tenantId string, eventId string, err error) error {
	return c.mongodb.eventService.UpdatePublishStateError(context.Background(), tenantId, eventId, err)
}

func (c *Service) publishMessage(req *es.ApplyEventRequest) error {
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
	return c.getPubsubAdapter().Publish(pubData)
}

func (c *Service) saveDomainEvent(req *es.ApplyEventRequest) error {
	event := &model.EventEntity{
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
	err := c.mongodb.eventService.Save(context.Background(), event)
	return err
}

func (c *Service) findEventById(tenantId string, id string) (*model.EventEntity, error) {
	return c.mongodb.eventService.FindById(context.Background(), tenantId, id)
}

func newError(msgType string, err error) error {
	return errors.New(msgType + err.Error())
}
