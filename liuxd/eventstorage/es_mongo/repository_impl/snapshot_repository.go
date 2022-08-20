package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
)

type snapshotRepository struct {
	dao *Dao[*model.Snapshot]
}

func NewSnapshotRepository(mongodb *db.MongoDbConfig, collName string) repository.SnapshotRepository {
	return &snapshotRepository{
		dao: NewDao[*model.Snapshot](mongodb, collName),
	}
}

func (r *snapshotRepository) Create(ctx context.Context, tenantId string, v *model.Snapshot) error {
	return r.dao.Insert(ctx, v)
}

func (r *snapshotRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *snapshotRepository) Update(ctx context.Context, tenantId string, v *model.Snapshot) error {
	return r.dao.Update(ctx, v)
}

func (r *snapshotRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Snapshot, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *snapshotRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) ([]*model.Snapshot, bool, error) {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
	}
	return r.dao.findList(ctx, tenantId, filter)
}

func (r *snapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.Snapshot, bool, error) {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	return r.dao.findOne(ctx, tenantId, filter)
}
