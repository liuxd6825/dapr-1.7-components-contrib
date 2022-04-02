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
	// LoadEvents 加载事件
	LoadEvents(ctx context.Context, req *LoadEventRequest) (*LoadResponse, error)
	// ExistAggregate 是否存在聚合根
	ExistAggregate(ctx context.Context, req *ExistAggregateRequest) (*ExistAggregateResponse, error)
	// ApplyEvent 应用事件
	ApplyEvent(ctx context.Context, req *ApplyEventRequest) (*ApplyResponse, error)
	// SaveSnapshot 保存镜像对象
	SaveSnapshot(ctx context.Context, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)
}
