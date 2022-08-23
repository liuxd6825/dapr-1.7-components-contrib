package dto

type PublishData struct {
	EventId        string      `json:"eventId"`
	EventData      interface{} `json:"eventData"`
	EventType      string      `json:"eventType"`
	EventVersion   string      `json:"eventVersion"`
	SequenceNumber uint64      `json:"sequenceNumber"`
}
