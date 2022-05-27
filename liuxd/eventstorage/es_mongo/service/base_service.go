package service

import "go.mongodb.org/mongo-driver/mongo"

const (
	IdField             = "_id"
	TenantIdField       = "tenant_id"
	AggregateIdField    = "aggregate_id"
	AggregateTypeField  = "aggregate_type"
	EventIdField        = "event_id"
	SequenceNumberField = "sequence_number"
	PublishStatusField  = "publish_status"
)

type BaseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}
