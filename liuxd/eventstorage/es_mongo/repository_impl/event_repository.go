package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository struct {
	dao *Dao[*model.Event]
}

func NewEventRepository(mongodb *db.MongoDbConfig, collName string) repository.EventRepository {
	return &eventRepository{
		dao: NewDao[*model.Event](mongodb, collName),
	}
}

func (r eventRepository) Create(ctx context.Context, tenantId string, v *model.Event) error {
	return r.dao.Insert(ctx, v)
}

func (r eventRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r eventRepository) Update(ctx context.Context, tenantId string, v *model.Event) error {
	return r.dao.Update(ctx, v)
}

func (r eventRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *eventRepository) UpdatePublishStatue(ctx context.Context, tenantId string, eventId string, publishStatue eventstorage.PublishStatus) error {
	filter := bson.M{IdField: eventId, TenantIdField: tenantId}
	data := bson.D{{"$set", bson.M{PublishStatusField: publishStatue}}}
	_, err := r.dao.getCollection(tenantId).UpdateOne(ctx, filter, data, options.Update())
	return err
}

func (r *eventRepository) FindByEventId(ctx context.Context, tenantId string, eventId string) (*model.Event, bool, error) {
	filter := bson.M{
		TenantIdField: tenantId,
		EventIdField:  eventId,
	}
	return r.dao.findOne(ctx, tenantId, filter)
}

func (r *eventRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	return r.dao.findList(ctx, tenantId, filter)
}

//
// FindNotPublishStatusSuccess
// @Description: 查找发送状态不成功的事件
// @receiver r
// @param ctx
// @param tenantId
// @param aggregateId
// @return *[]Event
// @return error
//
func (r *eventRepository) FindNotPublishStatusSuccess(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
		PublishStatusField: bson.M{"$ne": eventstorage.PublishStatusSuccess},
	}
	return r.dao.findList(ctx, tenantId, filter)
}

func (r *eventRepository) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error) {
	filter := bson.M{
		TenantIdField:       tenantId,
		AggregateIdField:    aggregateId,
		AggregateTypeField:  aggregateType,
		SequenceNumberField: bson.M{"$gt": sequenceNumber},
	}
	findOptions := options.Find().SetSort(bson.D{{SequenceNumberField, 1}})
	return r.dao.findList(ctx, tenantId, filter, findOptions)
}

/*
func (r *eventRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) uint64 {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	findOptions := options.FindOne().SetSort(bson.D{{SequenceNumberField, -1}})
	result := r.collection.FindOne(ctx, filter, findOptions)
	var event model.Event
	if err := result.Decode(&event); err == nil {
		return event.SequenceNumber + 1
	}
	return 1
}
*/
