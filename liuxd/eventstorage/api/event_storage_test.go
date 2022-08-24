package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl"
	pubsub_adapter "github.com/liuxd6825/components-contrib/pubsub"
	"github.com/liuxd6825/dapr/pkg/runtime/pubsub"
	"testing"
)

func TestEventStorage_FindEvents(t *testing.T) {
	es, err := newMongoEventSourcing()
	if err != nil {
		t.Error(err)
		return
	}
	req := &dto.FindEventsRequest{
		TenantId:      "001",
		AggregateType: "",
		Filter:        fmt.Sprintf(`aggregate_id=="%v"`, "inv-20200820-0001"),
		PageNum:       0,
		PageSize:      10,
		IsTotalRows:   false,
	}
	resp, err := es.FindEvents(context.Background(), req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("FindEvents() data.length=%v", len(resp.Data))
}

func TestEventStorage_DeleteEvent(t *testing.T) {
	es, err := newMongoEventSourcing()
	if err != nil {
		t.Error(err)
		return
	}
	createReq := newCreateEventRequest()
	if _, err := es.CreateEvent(context.Background(), createReq); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("CreateEvent() aggregateId=%v, tenantId=%v", createReq.AggregateId, createReq.TenantId)
	}

	applyReq := newApplyEventRequest(createReq.AggregateId)
	if _, err := es.ApplyEvent(context.Background(), applyReq); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("ApplyEvent() aggregateId=%v, tenantId=%v", applyReq.AggregateId, applyReq.TenantId)
	}

	deleteReq := newDeleteEventRequest(createReq.AggregateId)
	if _, err := es.DeleteEvent(context.Background(), deleteReq); err != nil {
		t.Error(err)
		return
	} else {
		t.Logf("DeleteEvent() aggregateId=%v, tenantId=%v", deleteReq.AggregateId, deleteReq.TenantId)
	}

	if resp, err := es.LoadEvent(context.Background(), newLoadEventRequest(createReq.AggregateId)); err != nil {
		t.Error(err)
	} else {
		last := len(*resp.Events)
		if last == 0 {
			t.Error("events.length = 0")
			return
		}
		events := *resp.Events
		sequenceNumber := events[last-1].SequenceNumber
		if sequenceNumber != uint64(last) {
			t.Error(fmt.Errorf("events.SequenceNumber %v != %v", sequenceNumber, last))
			return
		}
		t.Logf("LoadEvent() events=%v, aggregateId=%v, tenantId=%v", len(*resp.Events), createReq.AggregateId, createReq.TenantId)
	}

}

func newLoadEventRequest(aggregateId string) *dto.LoadEventRequest {
	return &dto.LoadEventRequest{
		TenantId:      "001",
		AggregateId:   aggregateId,
		AggregateType: "TEST",
	}
}

func newApplyEventRequest(aggregateId string) *dto.ApplyEventsRequest {
	return &dto.ApplyEventsRequest{
		TenantId:      "001",
		AggregateId:   aggregateId,
		AggregateType: "TEST",
		Events: []*dto.EventDto{
			&dto.EventDto{
				Metadata:     map[string]string{},
				CommandId:    uuid.NewString(),
				EventId:      uuid.NewString(),
				EventData:    map[string]interface{}{},
				EventType:    "TEST",
				EventVersion: "V1.0",
				PubsubName:   "pubsub",
				Topic:        "test_apply_event",
			},
		},
	}
}

func newCreateEventRequest() *dto.CreateEventRequest {
	return &dto.CreateEventRequest{
		TenantId:      "001",
		AggregateId:   uuid.New().String(),
		AggregateType: "TEST",
		Events: []*dto.EventDto{
			&dto.EventDto{
				Metadata:     map[string]string{},
				CommandId:    uuid.NewString(),
				EventId:      uuid.NewString(),
				EventData:    map[string]interface{}{},
				EventType:    "TEST",
				EventVersion: "V1.0",
				PubsubName:   "pubsub",
				Topic:        "test_create_event",
			},
		},
	}
}

func newDeleteEventRequest(aggregateId string) *dto.DeleteEventRequest {
	return &dto.DeleteEventRequest{
		TenantId:      "001",
		AggregateId:   aggregateId,
		AggregateType: "TEST",
		Event: &dto.EventDto{
			Metadata:     map[string]string{},
			CommandId:    uuid.NewString(),
			EventId:      uuid.NewString(),
			EventData:    map[string]interface{}{},
			EventType:    "TEST",
			EventVersion: "V1.0",
			PubsubName:   "pubsub",
			Topic:        "test_delete_event",
		},
	}
}

func newMongoEventSourcing() (eventstorage.EventStorage, error) {
	metadata := common.Metadata{
		Properties: map[string]string{
			"host":            "192.168.64.8:27018,192.168.64.8:27019,192.168.64.8:27020",
			"replica-set":     "mongors",
			"databaseName":    "dapr_esdb",
			"username":        "dapr",
			"password":        "123456",
			"snapshotTrigger": "250",
		},
	}
	adapter := func() pubsub.Adapter {
		return newAdapter()
	}
	es := NewEventStorage(nil)
	opts, err := mongo_impl.NewOptions(nil, metadata, adapter)
	if err != nil {
		return nil, err
	}
	if err = es.Init(opts); err != nil {
		return nil, err
	}
	return es, err
}

type testPubsub struct {
}

func (t testPubsub) Init(metadata pubsub_adapter.Metadata) error {
	return nil
}

func (t testPubsub) Features() []pubsub_adapter.Feature {
	return []pubsub_adapter.Feature{}
}

func (t testPubsub) Publish(req *pubsub_adapter.PublishRequest) error {
	return nil
}

func (t testPubsub) Subscribe(req pubsub_adapter.SubscribeRequest, handler pubsub_adapter.Handler) error {
	return nil
}

func (t testPubsub) Close() error {
	return nil
}

func newPubsub() pubsub_adapter.PubSub {
	return &testPubsub{}
}

type testAdapter struct {
}

func (t testAdapter) GetPubSub(pubsubName string) pubsub_adapter.PubSub {
	return newPubsub()
}

func (t testAdapter) Publish(req *pubsub_adapter.PublishRequest) error {
	return nil
}

func newAdapter() pubsub.Adapter {
	return &testAdapter{}
}
