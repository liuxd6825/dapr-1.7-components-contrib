package service

import (
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"golang.org/x/net/context"
	"testing"
)

func TestAggregateService_Create(t *testing.T) {
	service, err := newAggregateService()
	if err != nil {
		t.Error(t)
		return
	}
	id := model.NewObjectID()
	agg := &model.AggregateEntity{
		Id:             id,
		TenantId:       "001",
		AggregateId:    string(id),
		AggregateType:  "type",
		SequenceNumber: 1,
	}

	if err = service.Create(context.Background(), agg); err != nil {
		t.Error(err)
		return
	}

	if agg, err := service.NextSequenceNumber(context.Background(), agg.TenantId, agg.AggregateId, 1); err != nil {
		t.Error()
		return
	} else {
		fmt.Printf("NextSequenceNumber() : %v", agg.SequenceNumber+1)
	}

	if agg, err := service.Delete(context.Background(), agg.TenantId, agg.AggregateId); err != nil {
		t.Error()
		return
	} else {
		fmt.Printf("Delete() : %v", agg.AggregateId)
	}

	if agg, err := service.FindById(context.Background(), agg.TenantId, agg.AggregateId); err != nil {
		t.Error()
		return
	} else {
		fmt.Printf("FindById(): id=%v; deleted=%v; sn=%v", agg.Id, agg.Deleted, agg.SequenceNumber)
	}
}

func newAggregateService() (AggregateService, error) {
	metadata := common.Metadata{
		Properties: map[string]string{
			"host":         "192.168.64.8:27018, 192.168.64.8:27019, 192.168.64.8:27020",
			"username":     "query-example",
			"password":     "123456",
			"databaseName": "query-example",
		},
	}
	mongodb := db.NewMongoDB(nil)
	if err := mongodb.Init(metadata); err != nil {
		return nil, err
	}
	client := mongodb.GetClient()
	coll := mongodb.NewCollection("dapr_aggregate_test")
	service := NewAggregateService(client, coll)
	return service, nil
}
