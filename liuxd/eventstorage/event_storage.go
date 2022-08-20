package eventstorage

import (
	"context"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	pubsub_adapter "github.com/liuxd6825/dapr/pkg/runtime/pubsub"
)

type Session interface {
	UseTransaction(context.Context, SessionFunc) error
}

type SessionFunc func(ctx context.Context) error

type GetPubsubAdapter func() pubsub_adapter.Adapter

type Options struct {
	Metadata       common.Metadata
	PubsubAdapter  GetPubsubAdapter
	EventRepos     interface{}
	SnapshotRepos  interface{}
	AggregateRepos interface{}
	RelationRepos  interface{}
	MessageRepos   interface{}
	SnapshotCount  uint64
	Session        Session
}

// EventStorage 领域事件存储接口
type EventStorage interface {
	// Init 初始化
	Init(opts *Options) error

	GetLogger() logger.Logger

	// LoadEvent 加载事件
	LoadEvent(ctx context.Context, req *LoadEventRequest) (*LoadResponse, error)

	// CreateEvent 创建聚合事件
	CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error)

	// DeleteEvent 删除聚合事件
	DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error)

	// ApplyEvent 应用事件
	ApplyEvent(ctx context.Context, req *ApplyEventsRequest) (*ApplyEventsResponse, error)

	// SaveSnapshot 保存镜像对象
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)

	// GetRelations 获取聚合根关系
	GetRelations(ctx context.Context, req *GetRelationsRequest) (*GetRelationsResponse, error)
}
