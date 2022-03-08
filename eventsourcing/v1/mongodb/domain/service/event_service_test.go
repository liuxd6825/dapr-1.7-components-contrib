package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

var testTenantId = "tenant1"

const (
	TestEventName     = "dapr_event"
	TestSnapshotName  = "dapr_snapshot"
	TestTenantId      = "tenant_1"
	TestAggregateType = "system.user"
)

func TestEventService_Save(t *testing.T) {
	ctx := context.Background()

	client, err := getMongoClient()
	if err != nil {
		assert.NoError(t, err)
	}
	collection := client.Database("dapr_esdb").Collection(TestEventName)

	event := newEvent()
	service := NewEventService(client, collection)
	err = service.Save(ctx, event)
	assert.NoError(t, err)
}

func TestEventService_UpdatePublishStateOk(t *testing.T) {
	ctx := context.Background()

	client, err := getMongoClient()
	if err != nil {
		assert.NoError(t, err)
	}
	collection := client.Database("dapr_esdb").Collection("domain_event")

	id := "eventId_16"

	service := NewEventService(client, collection)
	err = service.UpdatePublishStateOk(ctx, testTenantId, id)
	assert.NoError(t, err)

	err = service.UpdatePublishStateError(ctx, testTenantId, id, errors.New("test error"))
	assert.NoError(t, err)

	event, err := service.FindById(ctx, testTenantId, id)
	assert.NoError(t, err)
	fmt.Println(event)
}

func TestEventService_Find(t *testing.T) {
	client, err := getMongoClient()
	assert.NoError(t, err)

	repos := repository.NewEventRepository(client, getEventCollection(client))

	data, err1 := repos.FindById(context.Background(), testTenantId, "1")
	assert.NoError(t, err1)
	fmt.Println(data)

	list, err2 := repos.FindByAggregateId(context.Background(), testTenantId, "aggregate_id")
	assert.NoError(t, err2)
	fmt.Println(list)
}

func TestEventService_FindBySequenceNumber(t *testing.T) {
	client, err := getMongoClient()
	if err != nil {
		assert.NoError(t, err)
	}
	service := NewEventService(client, getEventCollection(client))

	event := newEvent()
	err = service.Save(context.Background(), event)
	assert.NoError(t, err)

	list, err := service.FindBySequenceNumber(context.Background(), testTenantId, event.AggregateId, 1)
	assert.NoError(t, err)
	println(list)
}

func getEventCollection(client *mongo.Client) *mongo.Collection {
	return client.Database("dapr_esdb").Collection(TestEventName)
}

func newEvent() *model.EventEntity {
	id, _ := uuid.NewUUID()
	code, _ := uuid.NewUUID()
	event := model.EventEntity{
		TenantId:      TestTenantId,
		Id:            id.String(),
		EventId:       id.String(),
		Metadata:      map[string]string{"token": "null"},
		EventData:     map[string]interface{}{"id": "id", "name": "liuxd", "code": code.String()},
		EventRevision: "1.0",
		EventType:     "CreateUserEvent",
		AggregateId:   "001",
		AggregateType: TestAggregateType,
	}
	return &event
}

func getMongoClient() (*mongo.Client, error) {
	// 设置mongoDB客户端连接信息
	param := fmt.Sprintf("mongodb://dapr:123456@192.168.64.4:27017/dapr_esdb")
	clientOptions := options.Client().ApplyURI(param)

	// 建立客户端连接
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client, err
}
