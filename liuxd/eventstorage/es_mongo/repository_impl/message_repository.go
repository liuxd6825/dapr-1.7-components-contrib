package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MessageDbId = "system"
)

var msgOptions = NewOptions().SetDbId(MessageDbId)

type messageRepository struct {
	dao *dao[*model.Message]
}

func NewMessageRepository(mongodb *db.MongoDbConfig, collName string) repository.MessageRepository {
	res := &messageRepository{
		dao: NewDao[*model.Message](mongodb, collName),
	}
	return res
}

func (m *messageRepository) Create(ctx context.Context, v *model.Message) error {
	return m.dao.Insert(ctx, v, msgOptions)
}

func (m *messageRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return m.dao.DeleteById(ctx, tenantId, id, msgOptions)
}

func (m *messageRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
	}
	return m.dao.deleteByFilter(ctx, tenantId, filter)
}

func (m *messageRepository) Update(ctx context.Context, v *model.Message) error {
	return m.dao.Update(ctx, v, msgOptions)
}

func (m *messageRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Message, bool, error) {
	return m.dao.FindById(ctx, tenantId, id, msgOptions)
}

func (m *messageRepository) FindAll(ctx context.Context, limit *int64) ([]*model.Message, bool, error) {
	filter := bson.M{}
	options := NewOptions().SetDbId(MessageDbId).SetSort(bson.D{{"create_time", 0}})
	return m.dao.findList(ctx, msgOptions.DbId, filter, limit, options)
}
