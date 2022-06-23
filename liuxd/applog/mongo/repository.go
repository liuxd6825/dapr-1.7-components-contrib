package mongo

import (
	"context"
	"fmt"
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
	filter := bson.D{{"_id", entity.Id}}

	data := bson.D{{"$set", bson.M{
		"tenantId":  entity.TenantId,
		"appId":     entity.AppId,
		"class":     entity.Class,
		"func":      entity.Func,
		"level":     entity.Level,
		"time":      entity.Time,
		"status":    entity.Status,
		"message":   entity.Message,
		"pubAppId":  entity.PubAppId,
		"eventId":   entity.EventId,
		"commandId": entity.CommandId,
	}}}
	_, err := r.collection.UpdateOne(ctx, filter, data, options.Update())
	return err
}

func (r *EventLogRepository) FindById(ctx context.Context, tenantId, id string) (*EventLog, error) {
	var result EventLog
	filter := bson.D{
		{"_id", id},
		{"tenantId", tenantId},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err = getError(err); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *EventLogRepository) FindBySubAppIdAndCommandId(ctx context.Context, tenantId string, appId string, commandId string) (*[]EventLog, error) {
	filter := bson.D{
		{"tenantId", tenantId},
		{"appId", appId},
		{"commandId", commandId},
	}
	var list []EventLog
	cursor, err := r.collection.Find(ctx, filter)
	defer func() { // 关闭
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return &list, nil
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
	data := bson.D{{"$set", bson.M{
		"tenantId": entity.TenantId,
		"appId":    entity.AppId,
		"class":    entity.Class,
		"func":     entity.Func,
		"level":    entity.Level,
		"time":     entity.Time,
		"status":   entity.Status,
		"message":  entity.Message,
	}}}
	_, err := r.collection.UpdateOne(ctx, filter, data, options.Update())
	return err
}

func (r *AppLogRepository) FindById(ctx context.Context, tenantId string, id string) (*AppLog, error) {
	var result AppLog
	filter := bson.D{
		{"_id", id},
		{"tenantId", tenantId},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *AppLogRepository) FindByEventId(ctx context.Context, tenantId string, eventId string) (*[]AppLog, error) {
	filter := bson.D{
		{"_id", id},
		{"tenantId", tenantId},
		{"eventId", eventId},
	}
	list := []AppLog{}
	cursor, err := r.collection.Find(ctx, filter)
	if err = getError(err); err != nil {
		return nil, err
	}
	defer func() { // 关闭
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func getError(err error) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return err
}
