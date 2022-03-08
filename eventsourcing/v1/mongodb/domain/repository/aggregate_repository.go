package repository

import (
	"context"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	var event model.EventEntity
	err := r.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		return false, err
	}
	if &event == nil {
		return false, nil
	}
	return true, nil
}
