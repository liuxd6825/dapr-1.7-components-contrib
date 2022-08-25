package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl/db"
	cmap "github.com/orcaman/concurrent-map"
	"go.mongodb.org/mongo-driver/bson"
)

var collections = cmap.New()

type relationRepository struct {
	dao *dao[*model.Relation]
}

func NewRelationRepository(mongodb *db.MongoDbConfig, collName string) repository.RelationRepository {
	res := &relationRepository{
		dao: NewDao[*model.Relation](mongodb, collName),
	}
	return res
}

func (r *relationRepository) Create(ctx context.Context, tenantId string, v *model.Relation) error {
	return r.dao.Insert(ctx, v)
}

func (r *relationRepository) CreateMany(ctx context.Context, tenantId string, vList []*model.Relation) error {
	return r.dao.InsertMany(ctx, tenantId, vList)
}

func (r *relationRepository) DeleteById(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *relationRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
	}
	return r.dao.deleteByFilter(ctx, tenantId, filter)
}

func (r *relationRepository) Update(ctx context.Context, tenantId string, v *model.Relation) error {
	return r.dao.Update(ctx, v)
}

func (r *relationRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *relationRepository) FindPaging(ctx context.Context, query dto.FindPagingQuery) (*dto.FindPagingResult[*model.Relation], bool, error) {
	return r.dao.findPaging(ctx, query).Result()
}
