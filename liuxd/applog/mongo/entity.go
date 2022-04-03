package mongo

import "time"

type EventLog struct {
	Id       string     `json:"id" bson:"_id"`
	TenantId string     `json:"tenantId" bson:"tenantId"`
	AppId    string     `json:"appId" bson:"appId"`
	Class    string     `json:"class" bson:"class"`
	Func     string     `json:"func" bson:"func"`
	Level    string     `json:"level" bson:"level"`
	Time     *time.Time `json:"time" bson:"time"`
	Status   bool       `json:"status" bson:"status"`
	Message  string     `json:"message" bson:"message"`

	PubAppId  string `json:"pubAppId" bson:"pubAppId"`
	EventId   string `json:"eventId" bson:"eventId"`
	CommandId string `json:"commandId" bson:"commandId"`
}

type AppLog struct {
	Id       string     `json:"id" bson:"_id"`
	TenantId string     `json:"tenantId" bson:"tenantId"`
	AppId    string     `json:"appId" bson:"appId"`
	Class    string     `json:"class" bson:"class"`
	Func     string     `json:"func" bson:"func"`
	Level    string     `json:"level" bson:"level"`
	Time     *time.Time `json:"time" bson:"time"`
	Status   bool       `json:"status" bson:"status"`
	Message  string     `json:"message" bson:"message"`
}
