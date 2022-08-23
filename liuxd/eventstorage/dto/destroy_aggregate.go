package dto

type DeleteAggregateRequest struct {
	TenantId      string `json:"tenant_id"`
	AggregateType string `json:"aggregate_type"`
	AggregateId   string `json:"aggregate_id"`
}

type DeleteAggregateResponse struct {
	Headers *ResponseHeaders `json:"headers"`
}
