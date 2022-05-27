package model

type AggregateEntity struct {
	Id             ObjectID `bson:"_id"`
	TenantId       string   `bson:"tenant_id" `
	AggregateId    string   `bson:"aggregate_id"`
	AggregateType  string   `bson:"aggregate_type"`
	SequenceNumber uint64   `bson:"sequence_number"`
	Deleted        bool     `bson:"deleted"`
}
