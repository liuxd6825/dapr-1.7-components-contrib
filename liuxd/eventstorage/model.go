package eventstorage

import "go.mongodb.org/mongo-driver/bson/primitive"

type AggregateEntity interface {
	GetId() string
	GetTenantId() string
	GetAggregateId() string
	GetAggregateType() string
	GetSequenceNumber() uint64
}

type PublishStatus int

const (
	PublishStatusWait PublishStatus = iota
	PublishStatusSuccess
	PublishStatusError
)

func (s PublishStatus) ToString() string {
	switch s {
	case PublishStatusError:
		return "PublishStatusError"
	case PublishStatusSuccess:
		return "PublishStatusSuccess"
	case PublishStatusWait:
		return "PublishStatusWait"
	}
	return "NONE"
}

type EventEntity interface {
	GetId() string
	GetTenantId() string
	GetCommandId() string
	GetEventId() string
	GetMetadata() map[string]string
	GetEventData() map[string]interface{}
	GetEventType() string
	GetEventRevision() string
	GetAggregateId() string
	GetAggregateType() string
	GetSequenceNumber() uint64
	GetTimeStamp() primitive.DateTime
	GetTopic() string
	GetPublishName() string
	GetPublishStatus() PublishStatus
}

type SnapshotEntity interface {
	GetId() string
	GetTenantId() string
	GetAggregateId() string
	GetAggregateType() string
	GetAggregateData() map[string]interface{}
	GetAggregateRevision() string
	GetSequenceNumber() uint64
	GetMetadata() map[string]string
	GetTimeStamp() primitive.DateTime
}
