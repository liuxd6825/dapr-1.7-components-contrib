package api

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

func NewSnapshotDto(snapshotEntity *model.Snapshot) *eventstorage.LoadResponseSnapshotDto {
	return &eventstorage.LoadResponseSnapshotDto{
		AggregateData:    snapshotEntity.AggregateData,
		AggregateVersion: snapshotEntity.AggregateVersion,
		SequenceNumber:   snapshotEntity.SequenceNumber,
		Metadata:         snapshotEntity.Metadata,
	}
}

func NewEventDto(event *model.Event) *eventstorage.LoadResponseEventDto {
	return &eventstorage.LoadResponseEventDto{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventVersion:   event.EventVersion,
		SequenceNumber: event.SequenceNumber,
	}
}

func NewLoadResponse(tenantId string, aggregateId string, aggregateType string, snapshotEntity *model.Snapshot, events []*model.Event) *eventstorage.LoadResponse {
	var snapshotDto *eventstorage.LoadResponseSnapshotDto
	if snapshotEntity != nil {
		snapshotDto = NewSnapshotDto(snapshotEntity)
	}

	var eventsDto []eventstorage.LoadResponseEventDto
	if events != nil {
		eventsDto = make([]eventstorage.LoadResponseEventDto, len(events))
		for i, event := range events {
			eventsDto[i] = *NewEventDto(event)
		}
	}
	return &eventstorage.LoadResponse{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Snapshot:      snapshotDto,
		Events:        &eventsDto,
	}
}
