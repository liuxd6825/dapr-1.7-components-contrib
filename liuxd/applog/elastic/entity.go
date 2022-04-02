package elastic

type EventLog struct {
	TenantId  string `json:"tenantId" `
	Id        string `json:"id"`
	PubAppId  string `json:"pubAppId" `
	SubAppId  string `json:"subAppId" `
	EventId   string `json:"eventId"  `
	CommandId string `json:"commandId" `
	Status    bool   `json:"status"  `
	Message   string `json:"errorMsg" `
}

type AppLog struct {
	Id        string `json:"id" `
	TenantId  string `json:"tenantId"  `
	PubAppId  string `json:"pubAppId" `
	SubAppId  string `json:"subAppId" `
	EventId   string `json:"eventId" `
	CommandId string `json:"commandId" `
	Status    bool   `json:"status" `
	Message   string `json:"errorMsg"`
}
