package dto

import (
	"time"
)

type Event struct {
	TenantId      string                 `json:"tenantId"`
	AggregateId   string                 `json:"aggregateId"`
	AggregateType string                 `json:"aggregateType"`
	CommandId     string                 `json:"commandId"`
	EventId       string                 `json:"eventId"`
	EventData     map[string]interface{} `json:"eventData"`
	EventType     string                 `json:"eventType"`
	EventVersion  string                 `json:"eventVersion"`
	PubsubName    string                 `json:"pubsubName"`
	Relations     map[string]string      `json:"relations"`
	Topic         string                 `json:"topic"`
	Metadata      map[string]string      `json:"metadata"`
}

type EventDto struct {
	Metadata     map[string]string      `json:"metadata"`
	CommandId    string                 `json:"commandId"`
	EventId      string                 `json:"eventId"`
	EventData    map[string]interface{} `json:"eventData"`
	EventType    string                 `json:"eventType"`
	EventVersion string                 `json:"eventVision"`
	Relations    map[string]string      `json:"relations"`
	EventTime    time.Time              `json:"eventTime"`
	PubsubName   string                 `json:"pubsubName"`
	Topic        string                 `json:"topic"`
}
