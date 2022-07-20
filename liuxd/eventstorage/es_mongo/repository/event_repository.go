package repository

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	BaseRepository[*model.EventEntity]
}

func NewEventRepository(mongodb *db.MongoDB, collection *mongo.Collection) *EventRepository {
	res := &EventRepository{}
	res.mongodb = mongodb
	res.collection = collection
	return res
}

func (r *EventRepository) Insert(ctx context.Context, entity *model.EventEntity) error {
	idValue, err := model.ObjectIDFromHex(entity.EventId)
	if err != nil {
		return err
	}
	entity.Id = idValue
	_, err = r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepository) UpdatePublishStatue(ctx context.Context, eventId string, publishStatue eventstorage.PublishStatus) error {
	idValue, err := model.ObjectIDFromHex(eventId)
	if err != nil {
		return err
	}
	filter := bson.D{{IdField, idValue}}
	data := bson.D{{"$set", bson.M{PublishStatusField: publishStatue}}}
	_, err = r.collection.UpdateOne(ctx, filter, data, options.Update())
	return err
}

func (r *EventRepository) FindById(ctx context.Context, tenantId string, eventId string) (*model.EventEntity, error) {
	idValue, err := model.ObjectIDFromHex(eventId)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{IdField, idValue}}
	return r.findOne(ctx, filter)
}

/*func (r *EventRepository) FindByEventId(ctx context.Context, tenantId string, eventId string) (*model.EventEntity, error) {
	filter := bson.M{
		TenantIdField: tenantId,
		EventIdField:  eventId,
	}
	return r.findOne(ctx, filter)
}
*/

func (r *EventRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*[]model.EventEntity, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	return r.findList(ctx, filter)
}

//
// FindNotPublishStatusSuccess
// @Description: 查找发送状态不成功的事件
// @receiver r
// @param ctx
// @param tenantId
// @param aggregateId
// @return *[]EventEntity
// @return error
//
func (r *EventRepository) FindNotPublishStatusSuccess(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*[]model.EventEntity, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
		PublishStatusField: bson.M{"$ne": eventstorage.PublishStatusSuccess},
	}
	return r.findList(ctx, filter)
}

func (r *EventRepository) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) (*[]model.EventEntity, error) {
	filter := bson.M{
		TenantIdField:       tenantId,
		AggregateIdField:    aggregateId,
		AggregateTypeField:  aggregateType,
		SequenceNumberField: bson.M{"$gt": sequenceNumber},
	}
	findOptions := options.Find().SetSort(bson.D{{SequenceNumberField, 1}})
	return r.findList(ctx, filter, findOptions)
}

func (r *EventRepository) findList(ctx context.Context, filter interface{}, findOptions ...*options.FindOptions) (*[]model.EventEntity, error) {
	cursor, err := r.collection.Find(ctx, filter, findOptions...)
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	var list []model.EventEntity
	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *EventRepository) findOne(ctx context.Context, filter interface{}) (*model.EventEntity, error) {
	var result model.EventEntity
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

/*
func (r *EventRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) uint64 {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	findOptions := options.FindOne().SetSort(bson.D{{SequenceNumberField, -1}})
	result := r.collection.FindOne(ctx, filter, findOptions)
	var event model.EventEntity
	if err := result.Decode(&event); err == nil {
		return event.SequenceNumber + 1
	}
	return 1
}
*/
