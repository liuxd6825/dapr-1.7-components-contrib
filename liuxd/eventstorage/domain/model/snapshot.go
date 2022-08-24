package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Snapshot struct {
	Id               string                 `bson:"_id" json:"id"  gorm:"primaryKey"`
	TenantId         string                 `bson:"tenant_id" json:"tenant_id"`
	AggregateId      string                 `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType    string                 `bson:"aggregate_type" json:"aggregate_type"`
	AggregateData    map[string]interface{} `bson:"aggregate_data" json:"aggregate_data"  gorm:"type:text;serializer:json"`
	AggregateVersion string                 `bson:"aggregate_version" json:"aggregate_version"`
	SequenceNumber   uint64                 `bson:"sequence_number" json:"sequence_number"`
	Metadata         map[string]string      `bson:"metadata" json:"metadata"  gorm:"type:text;serializer:json"`
	TimeStamp        primitive.DateTime     `bson:"time_stamp" json:"time_stamp"`
}

func (r *Snapshot) GetId() string {
	return r.Id
}

func (r *Snapshot) SetId(v string) {
	r.Id = v
}

func (r *Snapshot) GetTenantId() string {
	return r.TenantId
}
