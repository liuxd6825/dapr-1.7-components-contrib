package service

import (
	"context"
	"fmt"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestSnapshotRepository_Insert(t *testing.T) {
	client, err := getMongoClient()
	assert.NoError(t, err)

	snapshot := newSnapshot("001", 3)
	repos := repository.NewSnapshotRepository(client, getSnapshotCollection(client))
	err = repos.Insert(context.Background(), snapshot)
	assert.NoError(t, err)
}

func TestSnapshotRepository_Find(t *testing.T) {
	client, err := getMongoClient()
	assert.NoError(t, err)

	repos := repository.NewSnapshotRepository(client, getSnapshotCollection(client))

	data, err1 := repos.FindByAggregateId(context.Background(), testTenantId, "aggregate_id")
	assert.NoError(t, err1)
	fmt.Println(data)
}

func TestSnapshotService_Save(t *testing.T) {
	client, err := getMongoClient()
	assert.NoError(t, err)

	ctx := context.Background()
	service := NewSnapshotService(client, getSnapshotCollection(client))

	snapshot := newSnapshot("001", 3)
	err = service.Save(ctx, snapshot)
	assert.NoError(t, err)

	snapshot, err1 := service.FindByMaxSequenceNumber(context.Background(), TestTenantId, "001")
	assert.NoError(t, err1)
	fmt.Println(snapshot)
}

func TestSnapshotService_FindByMaxSequenceNumber(t *testing.T) {
	client, err := getMongoClient()
	assert.NoError(t, err)

	ctx := context.Background()
	service := NewSnapshotService(client, getSnapshotCollection(client))

	snapshot := newSnapshot("001", 1)
	err = service.Save(ctx, snapshot)
	assert.NoError(t, err)

	snapshot, err1 := service.FindByMaxSequenceNumber(context.Background(), testTenantId, "aggregate_id_1")
	assert.NoError(t, err1)
	fmt.Println(snapshot)
}

func getSnapshotCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("dapr_esdb").Collection(TestSnapshotName)
}

func newSnapshot(aggregateId string, sequenceNumber int64) *model.SnapshotEntity {
	event := model.SnapshotEntity{
		TenantId:          TestTenantId,
		Metadata:          map[string]interface{}{"token": "null"},
		AggregateData:     map[string]interface{}{"id": "id", "name": "liuxd"},
		AggregateRevision: "1.0",
		AggregateType:     TestAggregateType,
		AggregateId:       aggregateId,
		SequenceNumber:    sequenceNumber,
	}
	return &event
}
