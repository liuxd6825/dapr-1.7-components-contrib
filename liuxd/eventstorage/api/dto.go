package api

import (
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
)

func NewSnapshotDto(snapshotEntity *model.Snapshot) *dto.LoadResponseSnapshotDto {
	return &dto.LoadResponseSnapshotDto{
		AggregateData:    snapshotEntity.AggregateData,
		AggregateVersion: snapshotEntity.AggregateVersion,
		SequenceNumber:   snapshotEntity.SequenceNumber,
		Metadata:         snapshotEntity.Metadata,
	}
}

func NewLoadResponseEventDto(event *model.Event) *dto.LoadResponseEventDto {
	return &dto.LoadResponseEventDto{
		EventId:        event.EventId,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventVersion:   event.EventVersion,
		SequenceNumber: event.SequenceNumber,
	}
}

func NewEventDto(tenantId string, aggregateId string, aggregateType string, event *dto.EventDto) (*dto.Event, error) {
	res := &dto.Event{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		CommandId:     event.CommandId,
		EventId:       event.EventId,
		EventData:     event.EventData,
		EventType:     event.EventType,
		EventVersion:  event.EventVersion,
		PubsubName:    event.PubsubName,
		Topic:         event.Topic,
		Metadata:      event.Metadata,
		Relations:     event.Relations,
	}
	return res, nil
}

func NewEvent(tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64, dto *dto.EventDto) *model.Event {
	event := &model.Event{
		Id:             dto.EventId,
		TenantId:       tenantId,
		CommandId:      dto.CommandId,
		EventId:        dto.EventId,
		Metadata:       dto.Metadata,
		EventData:      dto.EventData,
		EventType:      dto.EventType,
		EventVersion:   dto.EventVersion,
		AggregateId:    aggregateId,
		AggregateType:  aggregateType,
		SequenceNumber: sequenceNumber,
		Topic:          dto.Topic,
		PubsubName:     dto.PubsubName,
	}
	return event
}

func NewLoadResponse(tenantId string, aggregateId string, aggregateType string, snapshotEntity *model.Snapshot, events []*model.Event) *dto.LoadResponse {
	var snapshotDto *dto.LoadResponseSnapshotDto
	if snapshotEntity != nil {
		snapshotDto = NewSnapshotDto(snapshotEntity)
	}

	var eventsDto []dto.LoadResponseEventDto
	if events != nil {
		eventsDto = make([]dto.LoadResponseEventDto, len(events))
		for i, event := range events {
			eventsDto[i] = *NewLoadResponseEventDto(event)
		}
	}
	return &dto.LoadResponse{
		TenantId:      tenantId,
		AggregateId:   aggregateId,
		AggregateType: aggregateType,
		Snapshot:      snapshotDto,
		Events:        &eventsDto,
	}
}

func NewFindEventsItem(event *model.Event) *dto.FindEventsItem {
	res := &dto.FindEventsItem{
		Id:             event.Id,
		TenantId:       event.TenantId,
		CommandId:      event.CommandId,
		EventId:        event.EventId,
		Metadata:       event.Metadata,
		EventData:      event.EventData,
		EventType:      event.EventType,
		EventVersion:   event.EventVersion,
		AggregateId:    event.AggregateId,
		AggregateType:  event.AggregateType,
		SequenceNumber: event.SequenceNumber,
		TimeStamp:      event.TimeStamp,
		Topic:          event.Topic,
		PubsubName:     event.PubsubName,
	}
	return res
}

func NewFindEventsItems(events []*model.Event) []*dto.FindEventsItem {
	var list []*dto.FindEventsItem
	for _, item := range events {
		e := NewFindEventsItem(item)
		list = append(list, e)
	}
	return list
}
