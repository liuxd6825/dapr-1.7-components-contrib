package model

type EventTypeEntity struct {
	Id            ObjectID               `bson:"_id"`
	AppId         string                 `bson:"app_id"`
	TenantId      string                 `bson:"tenant_id" `
	EventType     string                 `bson:"event_type"`
	AggregateType string                 `bson:"aggregate_type"`
	Metadata      map[string]interface{} `bson:"metadata"`
	Version       string                 `bson:"version"`
}
