package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
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

func (r *AggregateRepository) FindById(ctx context.Context, tenantId string, aggregateId string) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	return r.findOne(ctx, filter)
}

func (r *AggregateRepository) Insert(ctx context.Context, aggregate *model.AggregateEntity) error {
	_, err := r.collection.InsertOne(ctx, aggregate)
	if err != nil {
		return err
	}
	return nil
}

func (r *AggregateRepository) Delete(ctx context.Context, tenantId, aggregateId string) error {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	update := bson.M{
		"$set": bson.M{"deleted": true},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *AggregateRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, count uint64) (*model.AggregateEntity, uint64, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, 0, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	update := bson.M{
		"$inc": bson.M{SequenceNumberField: count},
	}
	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	err = result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, 0, errors.New(fmt.Sprintf("aggregate idValue %s does not exist", aggregateId))
	} else if err != nil {
		return nil, 0, err
	}
	var aggregate model.AggregateEntity
	if err := result.Decode(&aggregate); err != nil {
		return nil, 0, err
	}
	return &aggregate, aggregate.SequenceNumber + 1, nil
}

func (r *AggregateRepository) findOne(ctx context.Context, filter interface{}) (*model.AggregateEntity, error) {
	var result model.AggregateEntity
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
