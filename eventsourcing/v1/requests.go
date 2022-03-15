package v1

type LoadEventRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
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

type SaveSnapshotRequest struct {
	TenantId          string                 `json:"tenantId"`
	AggregateId       string                 `json:"aggregateId"`
	AggregateType     string                 `json:"aggregateType"`
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type ExistAggregateRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
}
