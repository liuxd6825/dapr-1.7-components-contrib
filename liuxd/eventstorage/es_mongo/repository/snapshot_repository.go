package repository

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/other"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapshotRepository struct {
	BaseRepository[*model.SnapshotEntity]
}

func NewSnapshotRepository(mongodb *other.MongoDB, collection *mongo.Collection) *SnapshotRepository {
	res := &SnapshotRepository{}
	res.mongodb = mongodb
	res.collection = collection
	return res
}

func (r *SnapshotRepository) Insert(ctx context.Context, snapshot *model.SnapshotEntity) error {
	_, err := r.collection.InsertOne(ctx, snapshot)
	return err
}

func (r *SnapshotRepository) Update(ctx context.Context, snapshot *model.SnapshotEntity) error {
	_, err := r.collection.UpdateByID(ctx, snapshot.Id, snapshot)
	return err
}

func (r *SnapshotRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error) {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
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

	var list []model.SnapshotEntity
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *SnapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.SnapshotEntity, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	findOptions := options.FindOne().SetSort(bson.D{{SequenceNumberField, -1}})
	var snapshot model.SnapshotEntity
	if err := r.collection.FindOne(ctx, filter, findOptions).Decode(&snapshot); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}
