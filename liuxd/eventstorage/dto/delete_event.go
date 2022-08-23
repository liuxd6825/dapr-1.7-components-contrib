package dto

type DeleteEventRequest struct {
	TenantId      string    `json:"tenantId"`
	AggregateId   string    `json:"aggregateId"`
	AggregateType string    `json:"aggregateType"`
	Event         *EventDto `json:"event"`
}

type DeleteEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}
