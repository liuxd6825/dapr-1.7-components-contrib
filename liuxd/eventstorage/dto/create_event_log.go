package dto

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
	Headers *ResponseHeaders `json:"headers"`
}
