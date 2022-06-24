package eventstorage

import (
	"encoding/json"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"time"
)

type LoadEventRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
}

type LoadResponse struct {
	TenantId    string                   `json:"tenantId"`
	AggregateId string                   `json:"aggregateId"`
	Snapshot    *LoadResponseSnapshotDto `json:"snapshot"`
	Events      *[]LoadResponseEventDto  `json:"events"`
}

type LoadResponseSnapshotDto struct {
	AggregateData    map[string]interface{} `json:"aggregateData"`
	AggregateVersion string                 `json:"aggregateVersion"`
	SequenceNumber   uint64                 `json:"sequenceNumber"`
	Metadata         map[string]string      `json:"metadata"`
}

type LoadResponseEventDto struct {
	EventId        string                 `json:"eventId"`
	EventData      map[string]interface{} `json:"eventData"`
	EventType      string                 `json:"eventType"`
	EventVersion   string                 `json:"eventRevision"`
	SequenceNumber uint64                 `json:"sequenceNumber"`
}

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

func NewEvent(tenantId string, aggregateId string, aggregateType string, event EventDto) (*Event, error) {
	res := &Event{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		CommandId:     event.CommandId,
		EventId:       event.EventId,
		EventData:     event.EventData,
		EventType:     event.EventType,
		EventVersion:  event.EventVersion,
		PubsubName:    event.PubsubName,
		Topic:         event.Topic,
		Metadata:      event.Metadata,
		Relations:     event.Relations,
	}
	return res, nil
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

type ApplyEventsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	Events        *[]EventDto
}

type ApplyEventsResponse struct {
}

type CreateEventRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	Events        *[]EventDto
}

type CreateEventResponse struct {
}

type DeleteEventRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	Event         *EventDto
}

type DeleteEventResponse struct {
}

type SaveSnapshotRequest struct {
	TenantId         string                 `json:"tenantId"`
	AggregateId      string                 `json:"aggregateId"`
	AggregateType    string                 `json:"aggregateType"`
	AggregateData    map[string]interface{} `json:"aggregateData"`
	AggregateVersion string                 `json:"aggregateVersion"`
	SequenceNumber   uint64                 `json:"sequenceNumber"`
	Metadata         map[string]string      `json:"metadata"`
}

type SaveSnapshotResponse struct {
}

type ExistAggregateRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
}

type ExistAggregateResponse struct {
	IsExist bool `json:"isExist"`
}

type CreateEventLogRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"errorMsg"`
}

type CreateEventLogResponse struct {
}

type UpdateEventLogRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"errorMsg"`
}

type UpdateEventLogResponse struct {
}

type GetEventLogByCommandIdRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	CommandId string `json:"commandId"`
}

type GetEventLogByCommandIdResponse struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"message"`
}

type GetRelationsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
}

func (g *GetRelationsRequest) GetTenantId() string {
	return g.TenantId
}

func (g *GetRelationsRequest) GetFilter() string {
	return g.Filter
}

func (g *GetRelationsRequest) GetSort() string {
	return g.Sort
}

func (g *GetRelationsRequest) GetPageNum() uint64 {
	return g.PageNum
}
func (g *GetRelationsRequest) GetPageSize() uint64 {
	return g.PageSize
}

type GetRelationsResponse struct {
	Data       []*Relation `json:"data"`
	TotalRows  uint64      `json:"totalRows"`
	TotalPages uint64      `json:"totalPages"`
	PageNum    uint64      `json:"pageNum"`
	PageSize   uint64      `json:"pageSize"`
	Filter     string      `json:"filter"`
	Sort       string      `json:"sort"`
	Error      string      `json:"error"`
	IsFound    bool        `json:"isFound"`
}

type Relation struct {
	Id          string            `json:"id"`
	TenantId    string            `json:"tenantId"`
	TableName   string            `json:"tableName"`
	AggregateId string            `json:"aggregateId"`
	IsDeleted   bool              `json:"isDeleted"`
	Items       map[string]string `json:"items"`
}

func (r *Relation) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})
	data["id"] = r.Id
	data["tenantId"] = r.TenantId
	data["tableName"] = r.TableName
	data["aggregateId"] = r.AggregateId
	data["isDeleted"] = r.IsDeleted
	for k, v := range r.Items {
		name := utils.AsJsonName(k)
		data[name] = v
	}
	return json.Marshal(data)
}
