package mongo

import "fmt"

type EventLog struct {
	Id        string `json:"id" bson:"_id"`
	TenantId  string `json:"tenantId" bson:"tenantId"`
	PubAppId  string `json:"pubAppId" bson:"pubAppId"`
	SubAppId  string `json:"subAppId" bson:"subAppId"`
	EventId   string `json:"eventId" bson:"eventId"`
	CommandId string `json:"commandId" bson:"commandId"`
	Status    bool   `json:"status" bson:"status"`
	Message   string `json:"errorMsg" bson:"errorMsg"`
}

func (l *EventLog) GetId() string {
	return GetEventLogId(l.TenantId, l.SubAppId, l.CommandId)
}

func GetEventLogId(tenantId string, subAppId string, commandId string) string {
	return fmt.Sprintf("%s_%s_%s", tenantId, subAppId, commandId)
}

type AppLog struct {
	Id        string `json:"id" bson:"_id"`
	TenantId  string `json:"tenantId" bson:"tenantId"`
	PubAppId  string `json:"pubAppId" bson:"pubAppId"`
	SubAppId  string `json:"subAppId" bson:"subAppId"`
	EventId   string `json:"eventId" bson:"eventId"`
	CommandId string `json:"commandId" bson:"commandId"`
	Status    bool   `json:"status" bson:"status"`
	Message   string `json:"errorMsg" bson:"errorMsg"`
}
