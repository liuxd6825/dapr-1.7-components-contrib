package applog

import (
	"context"
	"github.com/dapr/components-contrib/liuxd/common"
	pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
)

type GetPubsubAdapter func() pubsub_adapter.Adapter

type Logger interface {
	Init(metadata common.Metadata, getPubsubAdapter GetPubsubAdapter) error

	WriteEventLog(ctx context.Context, req *WriteEventLogRequest) (*WriteEventLogResponse, error)
	UpdateEventLog(ctx context.Context, req *UpdateEventLogRequest) (*UpdateEventLogResponse, error)
	GetEventLogByCommandId(ctx context.Context, req *GetEventLogByCommandIdRequest) (*GetEventLogByCommandIdResponse, error)

	WriteAppLog(ctx context.Context, req *WriteAppLogRequest) (*WriteAppLogResponse, error)
	UpdateAppLog(ctx context.Context, req *UpdateAppLogRequest) (*UpdateAppLogResponse, error)
	GetAppLogById(ctx context.Context, req *GetAppLogByIdRequest) (*GetAppLogByIdResponse, error)
}

type WriteEventLogRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"message"`
}

type WriteEventLogResponse struct {
}

type UpdateEventLogRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"errorMsg"`
}

type UpdateEventLogResponse struct {
}

// GetLogByCommandId

type GetEventLogByCommandIdRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	CommandId string `json:"commandId"`
}

type GetEventLogByCommandIdResponse struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	EventId   string `json:"eventId"`
	CommandId string `json:"commandId"`
	Status    bool   `json:"status"`
	Message   string `json:"message"`
}

//

type WriteAppLogRequest struct {
	TenantId  string `json:"tenantId"`
	AppId     string `json:"appId"`
	LogId     string `json:"logId"`
	Level     string `json:"level"`
	WriteTime string `json:"writeTime"`
	Message   string `json:"message"`
}

type WriteAppLogResponse struct {
}

type UpdateAppLogRequest struct {
	TenantId  string `json:"tenantId"`
	AppId     string `json:"appId"`
	LogId     string `json:"logId"`
	Level     string `json:"level"`
	WriteTime string `json:"writeTime"`
	Message   string `json:"message"`
}

type UpdateAppLogResponse struct {
}

// GetLogByCommandId

type GetAppLogByIdRequest struct {
	TenantId  string `json:"tenantId"`
	PubAppId  string `json:"pubAppId"`
	SubAppId  string `json:"subAppId"`
	CommandId string `json:"commandId"`
}

type GetAppLogByIdResponse struct {
	TenantId  string `json:"tenantId"`
	AppId     string `json:"appId"`
	LogId     string `json:"logId"`
	Level     string `json:"level"`
	WriteTime string `json:"writeTime"`
	Message   string `json:"message"`
}
