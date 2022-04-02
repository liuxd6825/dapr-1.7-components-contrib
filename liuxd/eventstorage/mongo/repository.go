package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type AggregateRepository struct {
	BaseRepository
}

func NewAggregateRepository(client *mongo.Client, collection *mongo.Collection) *AggregateRepository {
	return &AggregateRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *AggregateRepository) ExistAggregate(ctx context.Context, tenantId string, aggregateId string) (bool, error) {
	filter := bson.M{
		"tenant_id":    tenantId,
		"aggregate_id": aggregateId,
	}
	var event EventEntity
	err := r.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return false, err
	}
	if &event == nil {
		return false, nil
	}
	return true, nil
}

type EventRepository struct {
	BaseRepository
}

func NewEventRepository(client *mongo.Client, collection *mongo.Collection) *EventRepository {
	return &EventRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *EventRepository) Insert(ctx context.Context, entity *EventEntity) error {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventRepository) UpdatePublishState(ctx context.Context, tenantId string, id string, state int, time primitive.DateTime, err error) error {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}

	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", bson.M{"publish_state": state, "publish_time": time, "publish_error": errMessage}}}
	_, err = r.collection.UpdateOne(ctx, filter, data, options.Update())
	return err
}

func (r *EventRepository) FindById(ctx context.Context, tenantId string, id string) (*EventEntity, error) {
	var result EventEntity
	filter := bson.D{{"_id", id}}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *EventRepository) FindByEventId(ctx context.Context, tenantId string, eventId string) (*EventEntity, error) {
	var result EventEntity
	filter := bson.M{
		"tenant_id": tenantId,
		"event_id":  eventId,
	}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *EventRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]EventEntity, error) {
	filter := bson.M{
		"tenant_id":    tenantId,
		"aggregate_id": aggregateId,
	}
	cursor, _ := r.collection.Find(ctx, filter)

	defer func() { // 关闭
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	list := []EventEntity{}
	err := cursor.All(ctx, &list) // 当然也可以用   next
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *EventRepository) FindBySequenceNumber(ctx context.Context, tenantId string, aggregateId string, sequenceNumber int64) (*[]EventEntity, error) {
	filter := bson.M{
		"tenant_id":       tenantId,
		"aggregate_id":    aggregateId,
		"sequence_number": bson.M{"$gt": sequenceNumber},
	}
	findOptions := options.Find().SetSort(bson.D{{"sequence_number", 1}})
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	defer func() { // 关闭
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	list := []EventEntity{}
	err = cursor.All(ctx, &list) // 当然也可以用   next
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *EventRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) int64 {
	filter := bson.M{
		"aggregate_id":   aggregateId,
		"aggregate_type": aggregateType,
	}
	findOptions := options.FindOne().SetSort(bson.D{{"sequence_number", -1}})
	result := r.collection.FindOne(ctx, filter, findOptions)
	var event EventEntity
	if err := result.Decode(&event); err == nil {
		return event.SequenceNumber + 1
	}
	return 1
}

type SnapshotRepository struct {
	BaseRepository
}

func NewSnapshotRepository(client *mongo.Client, collection *mongo.Collection) *SnapshotRepository {
	return &SnapshotRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *SnapshotRepository) Insert(ctx context.Context, snapshot *SnapshotEntity) error {
	_, err := r.collection.InsertOne(ctx, snapshot)
	if err != nil {
		return err
	}
	return nil
}

func (r *SnapshotRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]SnapshotEntity, error) {
	filter := bson.M{
		"tenant_id":    tenantId,
		"aggregate_id": aggregateId,
	}

	cursor, err := r.collection.Find(ctx, filter)
	defer func() { // 关闭
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	list := []SnapshotEntity{}
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *SnapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*SnapshotEntity, error) {
	filter := bson.M{
		"tenant_id":    tenantId,
		"aggregate_id": aggregateId,
	}
	findOptions := options.FindOne().SetSort(bson.D{{"sequence_number", -1}})
	var snapshot SnapshotEntity
	if err := r.collection.FindOne(ctx, filter, findOptions).Decode(&snapshot); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}
