package model

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Id             string                     `bson:"_id" json:"id"`
	TenantId       string                     `bson:"tenant_id" json:"tenant_id"`
	CommandId      string                     `bson:"command_id" json:"command_id"`
	EventId        string                     `bson:"event_id" json:"event_id"`
	Metadata       map[string]string          `bson:"metadata" json:"metadata"`
	EventData      map[string]interface{}     `bson:"event_data" json:"event_data"`
	EventType      string                     `bson:"event_type" json:"event_type"`
	EventVersion   string                     `bson:"event_version" json:"event_version"`
	AggregateId    string                     `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType  string                     `bson:"aggregate_type" json:"aggregate_type"`
	SequenceNumber uint64                     `bson:"sequence_number" json:"sequence_number"`
	TimeStamp      primitive.DateTime         `bson:"time_stamp" json:"time_stamp"`
	Topic          string                     `bson:"topic" json:"topic"`
	PublishName    string                     `bson:"publish_name" json:"publish_name"`
	PublishStatus  eventstorage.PublishStatus `bson:"publish_status" json:"publish_status"`
}

func (a *Event) GetId() string {
	return a.Id
}

func (a *Event) GetTenantId() string {
	return a.TenantId
}
