package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepository struct {
	dao *Dao[*model.Message]
}

func NewMessageRepository(mongodb *db.MongoDbConfig, collName string) repository.MessageRepository {
	res := &messageRepository{
		dao: NewDao[*model.Message](mongodb, collName),
	}
	return res
}

func (m messageRepository) Create(ctx context.Context, tenantId string, v *model.Message) error {
	return m.dao.Insert(ctx, v)
}

func (m messageRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return m.dao.DeleteById(ctx, tenantId, id)
}

func (m messageRepository) Update(ctx context.Context, tenantId string, v *model.Message) error {
	return m.dao.Update(ctx, v)
}

func (m messageRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Message, bool, error) {
	return m.dao.FindById(ctx, tenantId, id)
}

func (m messageRepository) FindSendList(ctx context.Context, tenantId string, maxResult int64) ([]*model.Message, bool, error) {
	filter := bson.M{
		TenantIdField: tenantId,
	}
	return m.dao.findList(ctx, tenantId, filter, options.Find().SetLimit(maxResult))
}
