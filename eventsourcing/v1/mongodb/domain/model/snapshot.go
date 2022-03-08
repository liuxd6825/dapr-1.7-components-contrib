package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SnapshotEntity struct {
	TenantId          string                 `bson:"tenant_id"`
	AggregateId       string                 `bson:"aggregate_id"`
	AggregateType     string                 `bson:"aggregate_type"`
	AggregateData     map[string]interface{} `bson:"aggregate_data"`
	AggregateRevision string                 `bson:"aggregate_revision"`
	SequenceNumber    int64                  `bson:"sequence_number"`
	Metadata          map[string]interface{} `bson:"metadata"`
	TimeStamp         primitive.DateTime     `bson:"time_stamp"`
}
