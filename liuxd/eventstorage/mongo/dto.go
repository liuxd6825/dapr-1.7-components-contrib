package mongo

import "github.com/dapr/components-contrib/liuxd/eventstorage"

func NewSnapshotDto(snapshotEntity *SnapshotEntity) *eventstorage.SnapshotDto {
	return &eventstorage.SnapshotDto{
		AggregateData:     snapshotEntity.AggregateData,
		AggregateRevision: snapshotEntity.AggregateRevision,
		SequenceNumber:    snapshotEntity.SequenceNumber,
		Metadata:          snapshotEntity.Metadata,
	}
}

func NewEventDto(event *EventEntity) *eventstorage.EventDto {
	return &eventstorage.EventDto{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventRevision:  event.EventRevision,
		SequenceNumber: event.SequenceNumber,
	}
}

func NewLoadResponse(tenantId string, aggregateId string, snapshotEntity *SnapshotEntity, events *[]EventEntity) *eventstorage.LoadResponse {
	var snapshotDto *eventstorage.SnapshotDto = nil
	if snapshotEntity != nil {
		snapshotDto = NewSnapshotDto(snapshotEntity)
	}

	var eventDtos []eventstorage.EventDto = nil
	if events != nil {
		eventDtos = make([]eventstorage.EventDto, len(*events))
		for i, event := range *events {
			eventDtos[i] = *NewEventDto(&event)
		}
	}
	return &eventstorage.LoadResponse{
		TenantId:    tenantId,
		AggregateId: aggregateId,
		Snapshot:    snapshotDto,
		Events:      &eventDtos,
	}
}
