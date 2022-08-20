package model

type EventType struct {
	Id            string                 `bson:"_id" json:"id"`
	AppId         string                 `bson:"app_id" json:"app_id"`
	TenantId      string                 `bson:"tenant_id" json:"tenant_id"`
	EventType     string                 `bson:"event_type" json:"event_type"`
	AggregateType string                 `bson:"aggregate_type" json:"aggregate_type"`
	Metadata      map[string]interface{} `bson:"metadata" json:"metadata"`
	Version       string                 `bson:"version" json:"version"`
}

func (a *EventType) GetId() string {
	return a.Id
}

func (a *EventType) GetTenantId() string {
	return a.TenantId
}
