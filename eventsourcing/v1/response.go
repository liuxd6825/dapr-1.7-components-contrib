package v1

import (
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
)

type LoadResponse struct {
	TenantId    string       `json:"tenantId"`
	AggregateId string       `json:"aggregateId"`
	Snapshot    *SnapshotDto `json:"snapshot"`
	Events      *[]EventDto  `json:"events"`
}

func NewLoadResponse(tenantId string, aggregateId string, snapshotEntity *model.SnapshotEntity, events *[]model.EventEntity) *LoadResponse {
	var snapshotDto *SnapshotDto = nil
	if snapshotEntity != nil {
		snapshotDto = NewSnapshotDto(snapshotEntity)
	}

	var eventDtos []EventDto = nil
	if events != nil {
		eventDtos = make([]EventDto, len(*events))
		for i, event := range *events {
			eventDtos[i] = *NewEventDto(&event)
		}
	}
	return &LoadResponse{
		TenantId:    tenantId,
		AggregateId: aggregateId,
		Snapshot:    snapshotDto,
		Events:      &eventDtos,
	}
}

type SnapshotDto struct {
	AggregateData     map[string]interface{} `json:"aggregateData"`
	AggregateRevision string                 `json:"aggregateRevision"`
	SequenceNumber    int64                  `json:"sequenceNumber"`
	Metadata          map[string]interface{} `json:"metadata"`
}

func NewSnapshotDto(snapshotEntity *model.SnapshotEntity) *SnapshotDto {
	return &SnapshotDto{
		AggregateData:     snapshotEntity.AggregateData,
		AggregateRevision: snapshotEntity.AggregateRevision,
		SequenceNumber:    snapshotEntity.SequenceNumber,
		Metadata:          snapshotEntity.Metadata,
	}
}

type EventDto struct {
	EventId        string                 `json:"eventId"`
	EventData      map[string]interface{} `json:"eventData"`
	EventType      string                 `json:"eventType"`
	EventRevision  string                 `json:"eventRevision"`
	SequenceNumber int64                  `json:"sequenceNumber"`
}

func NewEventDto(event *model.EventEntity) *EventDto {
	return &EventDto{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventRevision:  event.EventRevision,
		SequenceNumber: event.SequenceNumber,
	}
}

type ApplyResponse struct {
}

type SaveSnapshotResponse struct {
}

type ExistAggregateResponse struct {
	IsExist bool `json:"isExist"`
}
