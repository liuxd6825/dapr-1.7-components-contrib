package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindEventsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
	IsTotalRows   bool   `json:"isTotalRows"`
}

func (g *FindEventsRequest) GetTenantId() string {
	return g.TenantId
}

func (g *FindEventsRequest) GetFilter() string {
	return g.Filter
}

func (g *FindEventsRequest) GetSort() string {
	return g.Sort
}

func (g *FindEventsRequest) GetPageNum() uint64 {
	return g.PageNum
}

func (g *FindEventsRequest) GetPageSize() uint64 {
	return g.PageSize
}

func (g *FindEventsRequest) GetIsTotalRows() bool {
	return g.IsTotalRows
}

type FindEventsResponse struct {
	Data       []*FindEventsItem `json:"data"`
	Headers    *ResponseHeaders  `json:"headers"`
	TotalRows  uint64            `json:"totalRows"`
	TotalPages uint64            `json:"totalPages"`
	PageNum    uint64            `json:"pageNum"`
	PageSize   uint64            `json:"pageSize"`
	Filter     string            `json:"filter"`
	Sort       string            `json:"sort"`
	Error      string            `json:"error"`
	IsFound    bool              `json:"isFound"`
}

type FindEventsItem struct {
	Id             string                 `bson:"_id" json:"id"`
	TenantId       string                 `bson:"tenant_id" json:"tenant_id"`
	CommandId      string                 `bson:"command_id" json:"command_id"`
	EventId        string                 `bson:"event_id" json:"event_id"`
	Metadata       map[string]string      `bson:"metadata" json:"metadata"`
	EventData      map[string]interface{} `bson:"event_data" json:"event_data"`
	EventType      string                 `bson:"event_type" json:"event_type"`
	EventVersion   string                 `bson:"event_version" json:"event_version"`
	AggregateId    string                 `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType  string                 `bson:"aggregate_type" json:"aggregate_type"`
	SequenceNumber uint64                 `bson:"sequence_number" json:"sequence_number"`
	TimeStamp      primitive.DateTime     `bson:"time_stamp" json:"time_stamp"`
	Topic          string                 `bson:"topic" json:"topic"`
	PubsubName     string                 `bson:"pubsub_name" json:"pubsub_name"`
}
