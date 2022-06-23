package model

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventEntity struct {
	Id             string                     `bson:"_id"`
	TenantId       string                     `bson:"tenant_id" `
	CommandId      string                     `bson:"command_id"`
	EventId        string                     `bson:"event_id"`
	Metadata       map[string]string          `bson:"meta_data"`
	EventData      map[string]interface{}     `bson:"event_data"`
	EventType      string                     `bson:"event_type"`
	EventVersion   string                     `bson:"event_version"`
	AggregateId    string                     `bson:"aggregate_id"`
	AggregateType  string                     `bson:"aggregate_type"`
	SequenceNumber uint64                     `bson:"sequence_number"`
	TimeStamp      primitive.DateTime         `bson:"time_stamp"`
	Topic          string                     `bson:"topic"`
	PublishName    string                     `bson:"publish_name"`
	PublishStatus  eventstorage.PublishStatus `bson:"publish_status"`
}
