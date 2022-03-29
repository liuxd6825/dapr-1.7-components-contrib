package repository

import (
	"context"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (r *EventLogRepository) Insert(ctx context.Context, entity *model.EventLog) error {
	_, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventLogRepository) Update(ctx context.Context, entity *model.EventLog) error {
	id := entity.GetId()
	filter := bson.D{{"_id", id}}
	_, err := r.collection.UpdateOne(ctx, filter, entity, options.Update())
	return err
}

func (r *EventLogRepository) FindById(ctx context.Context, tenantId string, subAppId string, commandId string) (*model.EventLog, error) {
	var result model.EventLog
	id := model.GetEventLogId(tenantId, subAppId, commandId)
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
