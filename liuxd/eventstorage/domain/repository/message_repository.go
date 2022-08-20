package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type MessageRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Message) error
	Delete(ctx context.Context, tenantId string, id string) error
	Update(ctx context.Context, tenantId string, v *model.Message) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Message, bool, error)
	FindSendList(ctx context.Context, tenantId string, maxResult int64) ([]*model.Message, bool, error)
}
