package repository

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (r *SnapshotRepository) Insert(ctx context.Context, snapshot *model.SnapshotEntity) error {
	_, err := r.collection.InsertOne(ctx, snapshot)
	if err != nil {
		return err
	}
	return nil
}

func (r *SnapshotRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error) {
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

	list := []model.SnapshotEntity{}
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *SnapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*model.SnapshotEntity, error) {
	filter := bson.M{
		"tenant_id":    tenantId,
		"aggregate_id": aggregateId,
	}
	findOptions := options.FindOne().SetSort(bson.D{{"sequence_number", -1}})
	var snapshot model.SnapshotEntity
	if err := r.collection.FindOne(ctx, filter, findOptions).Decode(&snapshot); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &snapshot, nil
}
