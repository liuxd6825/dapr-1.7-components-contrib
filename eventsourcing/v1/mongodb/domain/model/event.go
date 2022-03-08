package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventEntity struct {
	TenantId       string                 `bson:"tenant_id"`
	Id             string                 `bson:"_id"`
	CommandId      string                 `json:"commandId"`
	EventId        string                 `bson:"event_id"`
	Metadata       map[string]string      `bson:"meta_data"`
	EventData      map[string]interface{} `bson:"event_data"`
	EventType      string                 `bson:"event_type"`
	EventRevision  string                 `bson:"event_revision"`
	AggregateId    string                 `bson:"aggregate_id"`
	AggregateType  string                 `bson:"aggregate_type"`
	SequenceNumber int64                  `bson:"sequence_number"`
	TimeStamp      primitive.DateTime     `bson:"time_stamp"`
	PublishName    string                 `bson:"publish_name"`
	Topic          string                 `bson:"topic"`
}
