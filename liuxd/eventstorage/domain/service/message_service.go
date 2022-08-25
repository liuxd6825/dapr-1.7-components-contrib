package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type MessageService interface {
	Create(ctx context.Context, msg *model.Message) error
	Delete(ctx context.Context, tenantId, id string) error
	DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error
	FindAll(ctx context.Context, limit *int64) ([]*model.Message, bool, error)
}

func NewMessageService(repos repository.MessageRepository) MessageService {
	return &messageService{repos: repos}
}

type messageService struct {
	repos repository.MessageRepository
}

func (m *messageService) Create(ctx context.Context, msg *model.Message) error {
	if msg == nil {
		return nil
	}
	return m.repos.Create(ctx, msg)
}

func (m *messageService) Delete(ctx context.Context, tenantId, id string) error {
	return m.repos.DeleteById(ctx, tenantId, id)
}

func (m *messageService) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	return m.repos.DeleteByAggregateId(ctx, tenantId, aggregateId)
}

func (m *messageService) FindAll(ctx context.Context, limit *int64) ([]*model.Message, bool, error) {
	return m.repos.FindAll(ctx, limit)
}
