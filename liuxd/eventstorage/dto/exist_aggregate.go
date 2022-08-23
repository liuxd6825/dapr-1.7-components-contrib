package dto

type ExistAggregateRequest struct {
	TenantId    string `json:"tenantId"`
	AggregateId string `json:"aggregateId"`
}

type ExistAggregateResponse struct {
	Headers *ResponseHeaders `json:"headers"`
	IsExist bool             `json:"isExist"`
}
