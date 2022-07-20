package service

import (
	"github.com/google/uuid"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"golang.org/x/net/context"
	"testing"
)

func TestAggregateService_Create(t *testing.T) {
	mongodb, err := newTestMongoDb()
	if err != nil {
		t.Error(err)
		return
	}

	coll := mongodb.NewCollection("dapr_aggregate_test")
	service := NewAggregateService(mongodb, coll)
	agg := &model.AggregateEntity{
		Id:             newId(),
		TenantId:       "001",
		AggregateId:    newId(),
		AggregateType:  "type",
		SequenceNumber: 1,
	}
	err = service.Create(context.Background(), agg)
	if err != nil {
		t.Error(err)
	}
}

func newId() string {
	return uuid.New().String()
}

func newTestMongoDb() (*db.MongoDB, error) {
	metadata := common.Metadata{
		Properties: map[string]string{
			"host":         "192.168.64.8:27018 192.168.64.8:27019 192.168.64.8:27020",
			"username":     "query-example",
			"password":     "123456",
			"databaseName": "query-example",
		},
	}
	mongodb := db.NewMongoDB(nil)
	if err := mongodb.Init(metadata); err != nil {
		return nil, err
	}
	return mongodb, nil
}
