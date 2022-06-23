package service

import (
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/other"
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
	id := model.NewObjectID()
	agg := &model.AggregateEntity{
		Id:             id,
		TenantId:       "001",
		AggregateId:    id,
		AggregateType:  "type",
		SequenceNumber: 1,
	}
	err = service.Create(context.Background(), agg)
	if err != nil {
		t.Error(err)
	}
}

func TestAggregateService_FindById(t *testing.T) {

}

func TestAggregateService_NextSequenceNumber(t *testing.T) {

}

func TestAggregateService_ExistAggregate(t *testing.T) {

}

func newTestMongoDb() (*other.MongoDB, error) {
	metadata := common.Metadata{
		Properties: map[string]string{
			"host":         "192.168.64.8:27018,192.168.64.8:27019,192.168.64.8:27020",
			"username":     "query-example",
			"password":     "123456",
			"replicaSet":   "mongors",
			"databaseName": "query-example",
		},
	}
	mongodb := other.NewMongoDB(nil)
	if err := mongodb.Init(metadata); err != nil {
		return nil, err
	}
	return mongodb, nil
}
