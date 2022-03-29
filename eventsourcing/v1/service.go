package v1

import pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
import "github.com/valyala/fasthttp"

type GetPubsubAdapter func() pubsub_adapter.Adapter

type EventSourcingService interface {
	Init(metadata Metadata, getAdapter GetPubsubAdapter) error
	LoadEvents(reqCtx *fasthttp.RequestCtx, req *LoadEventRequest) (*LoadResponse, error)
	ExistAggregate(reqCtx *fasthttp.RequestCtx, req *ExistAggregateRequest) (*ExistAggregateResponse, error)
	ApplyEvent(reqCtx *fasthttp.RequestCtx, req *ApplyEventRequest) (*ApplyResponse, error)
	SaveSnapshot(reqCtx *fasthttp.RequestCtx, req *SaveSnapshotRequest) (*SaveSnapshotResponse, error)

	CreateEventLog(reqCtx *fasthttp.RequestCtx, req *CreateEventLogRequest) (*CreateEventLogResponse, error)
	UpdateEventLog(reqCtx *fasthttp.RequestCtx, req *UpdateEventLogRequest) (*UpdateEventLogResponse, error)
	GetEventLogByCommandId(reqCtx *fasthttp.RequestCtx, req *GetEventLogByCommandIdRequest) (*GetEventLogByCommandIdResponse, error)
}
