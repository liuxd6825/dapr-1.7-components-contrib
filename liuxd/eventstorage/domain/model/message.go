package model

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"time"
)

type Message struct {
	Id          string              `bson:"_id" json:"id"`
	AggregateId string              `bson:"aggregate_id" json:"aggregate_id"`
	TenantId    string              `bson:"tenant_id" json:"tenant_id"`
	EventId     string              `bson:"event_id" json:"event_id"`
	CreateTime  time.Time           `bson:"create_time" json:"create_time"`
	Event       *eventstorage.Event `bson:"event" json:"event"`
}

func (a *Message) GetId() string {
	return a.Id
}

func (a *Message) GetTenantId() string {
	return a.TenantId
}
