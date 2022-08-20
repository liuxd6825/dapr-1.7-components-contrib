package repository_impl

import (
	"github.com/google/uuid"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"golang.org/x/net/context"
	"testing"
)

const (
	TenantId = "test"
)

func TestAggregateService_Create(t *testing.T) {
	mongodb, err := newTestMongoDb()
	if err != nil {
		t.Error(err)
		return
	}

	repos := NewAggregateRepository(mongodb, "aggregate")
	agg := &model.Aggregate{
		Id:             newId(),
		TenantId:       "test",
		AggregateId:    newId(),
		AggregateType:  "type",
		SequenceNumber: 1,
	}
	err = repos.Create(context.Background(), TenantId, agg)
	if err != nil {
		t.Error(err)
	}
}

func newId() string {
	return uuid.New().String()
}

func newTestMongoDb() (*db.MongoDbConfig, error) {
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
