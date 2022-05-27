package eventstorage

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
	}
	return res, nil
}

type EventDto struct {
	Metadata     map[string]string      `json:"metadata"`
	CommandId    string                 `json:"commandId"`
	EventId      string                 `json:"eventId"`
	EventData    map[string]interface{} `json:"eventData"`
	EventType    string                 `json:"eventType"`
	EventVersion string                 `json:"eventRevision"`
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
