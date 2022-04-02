package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type EventLogRepository struct {
	BaseRepository
}

func NewEventLogRepository(client *mongo.Client, collection *mongo.Collection) *EventLogRepository {
	return &EventLogRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *EventLogRepository) Insert(ctx context.Context, entity *EventLog) error {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventLogRepository) Update(ctx context.Context, entity *EventLog) error {
	id := entity.GetId()
	filter := bson.D{{"_id", id}}
	_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
	return err
}

func (r *EventLogRepository) FindById(ctx context.Context, tenantId string, subAppId string, commandId string) (*EventLog, error) {
	var result EventLog
	id := GetEventLogId(tenantId, subAppId, commandId)
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

type AppLogRepository struct {
	BaseRepository
}

func NewAppLogRepository(client *mongo.Client, collection *mongo.Collection) *AppLogRepository {
	return &AppLogRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *AppLogRepository) Insert(ctx context.Context, entity *AppLog) error {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (r *AppLogRepository) Update(ctx context.Context, entity *AppLog) error {
	id := entity.Id
	filter := bson.D{{"_id", id}}
	_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
	return err
}

func (r *AppLogRepository) FindById(ctx context.Context, tenantId string, subAppId string, commandId string) (*AppLog, error) {
	var result AppLog
	id := GetEventLogId(tenantId, subAppId, commandId)
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
