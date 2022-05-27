package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SnapshotEntity struct {
	Id               ObjectID               `bson:"_id"`
	TenantId         string                 `bson:"tenant_id"`
	AggregateId      string                 `bson:"aggregate_id"`
	AggregateType    string                 `bson:"aggregate_type"`
	AggregateData    map[string]interface{} `bson:"aggregate_data"`
	AggregateVersion string                 `bson:"aggregate_version"`
	SequenceNumber   uint64                 `bson:"sequence_number"`
	Metadata         map[string]string      `bson:"metadata"`
	TimeStamp        primitive.DateTime     `bson:"time_stamp"`
}
