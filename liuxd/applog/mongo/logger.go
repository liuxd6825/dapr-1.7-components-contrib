package mongo

import (
	"context"
	"github.com/dapr/components-contrib/liuxd/applog"
	"github.com/dapr/components-contrib/liuxd/common"
	pubsub_adapter "github.com/dapr/dapr/pkg/runtime/pubsub"
	"github.com/dapr/kit/logger"
)

type Logger struct {
	eventLogService EventLogService
	appLogService   AppLogService
	pubsubAdapter   pubsub_adapter.Adapter
	metadata        common.Metadata
	mongodb         *MongoDB
	log             logger.Logger
}

func NewLogger(log logger.Logger) applog.Logger {
	return &Logger{
		log: log,
	}
}

func (l Logger) Init(metadata common.Metadata, getPubsubAdapter applog.GetPubsubAdapter) error {
	l.pubsubAdapter = getPubsubAdapter()
	l.mongodb = NewMongoDB(l.log)
	l.metadata = metadata
	if err := l.mongodb.Init(metadata); err != nil {
		return err
	}

	mongoClient := l.mongodb.GetClient()

	appLogCollection := l.mongodb.NewCollection(l.mongodb.loggerMetadata.appLogCollectionName)
	eventLogCollection := l.mongodb.NewCollection(l.mongodb.loggerMetadata.eventLogCollectionName)

	l.appLogService = NewAppLogService(mongoClient, appLogCollection)
	l.eventLogService = NewEventLogService(mongoClient, eventLogCollection)
	return nil
}

func (l Logger) WriteAppLog(ctx context.Context, req *applog.WriteAppLogRequest) (*applog.WriteAppLogResponse, error) {
	return nil, nil
}

func (l Logger) UpdateAppLog(ctx context.Context, req *applog.UpdateAppLogRequest) (*applog.UpdateAppLogResponse, error) {
	return nil, nil
}

func (l Logger) GetAppLogById(ctx context.Context, req *applog.GetAppLogByIdRequest) (*applog.GetAppLogByIdResponse, error) {
	return nil, nil
}

func (l *Logger) WriteEventLog(ctx context.Context, req *applog.WriteEventLogRequest) (*applog.WriteEventLogResponse, error) {
	log := &EventLog{
		TenantId:  req.TenantId,
		PubAppId:  req.PubAppId,
		SubAppId:  req.SubAppId,
		EventId:   req.EventId,
		CommandId: req.CommandId,
		Status:    req.Status,
		Message:   req.Message,
	}
	err := l.eventLogService.Insert(ctx, log)
	if err != nil {
		return nil, err
	}
	return &applog.WriteEventLogResponse{}, nil
}

func (l *Logger) UpdateEventLog(ctx context.Context, req *applog.UpdateEventLogRequest) (*applog.UpdateEventLogResponse, error) {
	log := &EventLog{
		TenantId:  req.TenantId,
		PubAppId:  req.PubAppId,
		SubAppId:  req.SubAppId,
		EventId:   req.EventId,
		CommandId: req.CommandId,
		Status:    req.Status,
		Message:   req.Message,
	}
	err := l.eventLogService.Update(ctx, log)
	if err != nil {
		return nil, err
	}
	return &applog.UpdateEventLogResponse{}, nil
}

func (l *Logger) GetEventLogByCommandId(ctx context.Context, req *applog.GetEventLogByCommandIdRequest) (*applog.GetEventLogByCommandIdResponse, error) {
	tenantId := req.TenantId
	subAppId := req.SubAppId
	commandId := req.CommandId
	eventLog, err := l.eventLogService.FindById(ctx, tenantId, subAppId, commandId)
	if err != nil {
		return nil, err
	}
	return &applog.GetEventLogByCommandIdResponse{
		TenantId:  eventLog.TenantId,
		SubAppId:  eventLog.SubAppId,
		PubAppId:  eventLog.PubAppId,
		EventId:   eventLog.EventId,
		CommandId: eventLog.CommandId,
		Status:    eventLog.Status,
		Message:   eventLog.Message,
	}, nil
}
