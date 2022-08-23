package dto

type CreateEventRequest struct {
	TenantId      string      `json:"tenantId"`
	AggregateId   string      `json:"aggregateId"`
	AggregateType string      `json:"aggregateType"`
	Events        []*EventDto `json:"events"`
}

type CreateEventResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}
