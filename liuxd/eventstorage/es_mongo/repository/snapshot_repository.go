package repository

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/eventstorage/es_mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IdField             = "_id"
	TenantIdField       = "tenant_id"
	AggregateIdField    = "aggregate_id"
	AggregateTypeField  = "aggregate_type"
	EventIdField        = "event_id"
	SequenceNumberField = "sequence_number"
	PublishStatusField  = "publish_status"
)

type BaseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
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

func (r *SnapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*model.SnapshotEntity, error) {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
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
