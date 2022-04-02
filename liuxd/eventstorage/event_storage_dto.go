package eventstorage

type LoadEventRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
}

type LoadResponse struct {
	TenantId    string       `json:"tenantId"`
	AggregateId string       `json:"aggregateId"`
	Snapshot    *SnapshotDto `json:"snapshot"`
	Events      *[]EventDto  `json:"events"`
}

type SnapshotDto struct {
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type EventDto struct {
	EventId        string                 `json:"eventId"`
	EventData      map[string]interface{} `json:"eventData"`
	EventType      string                 `json:"eventType"`
	EventRevision  string                 `json:"eventRevision"`
	SequenceNumber int64                  `json:"sequenceNumber"`
}

type ApplyEventRequest struct {
	TenantId      string                 `json:"tenantId"`
	Metadata      map[string]string      `json:"metadata"`
	CommandId     string                 `json:"commandId"`
	EventId       string                 `json:"eventId"`
	EventData     map[string]interface{} `json:"eventData"`
	EventType     string                 `json:"eventType"`
	EventRevision string                 `json:"eventRevision"`
	AggregateId   string                 `json:"aggregateId"`
	AggregateType string                 `json:"aggregateType"`
	PubsubName    string                 `json:"pubsubName"`
	Topic         string                 `json:"topic"`
}

type ApplyResponse struct {
}

type SaveSnapshotRequest struct {
	TenantId          string                 `json:"tenantId"`
	AggregateId       string                 `json:"aggregateId"`
	AggregateType     string                 `json:"aggregateType"`
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
	Metadata          map[string]interface{} `json:"metadata"`
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
