package repository_impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

const (
	MessageDbId = "system"
)

var msgOptions = NewOptions().SetDbId(MessageDbId)

type messageRepository struct {
	dao *dao[*model.Message]
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	_ = db.AutoMigrate(&model.Message{})
	res := &messageRepository{
		dao: NewDao[*model.Message](db,
			func() *model.Message { return &model.Message{} },
			func() []*model.Message { return []*model.Message{} },
		),
	}
	return res
}

func (m *messageRepository) Create(ctx context.Context, v *model.Message) error {
	return m.dao.Insert(ctx, v, msgOptions)
}

func (m *messageRepository) DeleteById(ctx context.Context, tenantId string, id string) error {
	return m.dao.DeleteById(ctx, tenantId, id, msgOptions)
}

func (m *messageRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	where := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v"`, tenantId, aggregateId)
	return m.dao.deleteByFilter(ctx, tenantId, where)
}

func (m *messageRepository) Update(ctx context.Context, v *model.Message) error {
	return m.dao.Update(ctx, v, msgOptions)
}

func (m *messageRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Message, bool, error) {
	return m.dao.FindById(ctx, tenantId, id, msgOptions)
}

func (m *messageRepository) FindAll(ctx context.Context, limit *int64) ([]*model.Message, bool, error) {
	options := NewOptions().SetDbId(MessageDbId).SetSort(bson.D{{"create_time", 0}})
	return m.dao.findList(ctx, msgOptions.DbId, nil, limit, options)
}
