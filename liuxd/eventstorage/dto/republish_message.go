package dto

type RepublishMessageRequest struct {
	Limit int64 `json:"limit"`
}

type RepublishMessageResponse struct {
	Count int64  `json:"count"`
	Error string `json:"error"`
}
