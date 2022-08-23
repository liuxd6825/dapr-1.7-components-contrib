package dto

type GetEventLogByCommandIdRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	CommandId string `json:"commandId"`
}

type GetEventLogByCommandIdResponse struct {
	Headers   *ResponseHeaders `json:"headers"`
	TenantId  string           `json:"tenantId"`
	PubAppId  string           `json:"pubAppId"`
	SubAppId  string           `json:"subAppId"`
	EventId   string           `json:"eventId"`
	CommandId string           `json:"commandId"`
	Status    bool             `json:"status"`
	Message   string           `json:"message"`
}
