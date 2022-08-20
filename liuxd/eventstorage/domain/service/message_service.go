package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type MessageService interface {
	Create(ctx context.Context, msg *model.Message) error
	Delete(ctx context.Context, tenantId, id string) error
	FindSentList(ctx context.Context, tenantId string, max int64) ([]*model.Message, bool, error)
}

func NewMessageService(repos repository.MessageRepository) MessageService {
	return &messageService{repos: repos}
}

type messageService struct {
	repos repository.MessageRepository
}

func (m *messageService) Create(ctx context.Context, msg *model.Message) error {
	return m.repos.Create(ctx, msg.TenantId, msg)
}

func (m *messageService) Delete(ctx context.Context, tenantId, id string) error {
	return m.repos.Delete(ctx, tenantId, id)
}

func (m *messageService) FindSentList(ctx context.Context, tenantId string, maxResult int64) ([]*model.Message, bool, error) {
	return m.repos.FindSendList(ctx, tenantId, maxResult)
}
