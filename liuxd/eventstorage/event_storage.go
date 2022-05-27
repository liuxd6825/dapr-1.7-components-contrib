package eventstorage

import (
	"context"
	"github.com/dapr/components-contrib/liuxd/common"
	pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
)

type GetPubsubAdapter func() pubsub_adapter.Adapter

// EventStorage 领域事件存储接口
type EventStorage interface {
	// Init 初始化
	Init(metadata common.Metadata, getAdapter GetPubsubAdapter) error

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
}
