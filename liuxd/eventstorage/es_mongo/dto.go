package es_mongo

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
)

func NewSnapshotDto(snapshotEntity *model.SnapshotEntity) *eventstorage.LoadResponseSnapshotDto {
	return &eventstorage.LoadResponseSnapshotDto{
		AggregateData:    snapshotEntity.AggregateData,
		AggregateVersion: snapshotEntity.AggregateVersion,
		SequenceNumber:   snapshotEntity.SequenceNumber,
		Metadata:         snapshotEntity.Metadata,
	}
}

func NewEventDto(event *model.EventEntity) *eventstorage.LoadResponseEventDto {
	return &eventstorage.LoadResponseEventDto{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventVersion:   event.EventVersion,
		SequenceNumber: event.SequenceNumber,
	}
}

func NewLoadResponse(tenantId string, aggregateId string, aggregateType string, snapshotEntity *model.SnapshotEntity, events *[]model.EventEntity) *eventstorage.LoadResponse {
	var snapshotDto *eventstorage.LoadResponseSnapshotDto = nil
	if snapshotEntity != nil {
		snapshotDto = NewSnapshotDto(snapshotEntity)
	}

	var eventDtos []eventstorage.LoadResponseEventDto = nil
	if events != nil {
		eventDtos = make([]eventstorage.LoadResponseEventDto, len(*events))
		for i, event := range *events {
			eventDtos[i] = *NewEventDto(&event)
		}
	}
	return &eventstorage.LoadResponse{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Snapshot:      snapshotDto,
		Events:        &eventDtos,
	}
}
