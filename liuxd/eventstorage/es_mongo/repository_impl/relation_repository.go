package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	cmap "github.com/orcaman/concurrent-map"
)

var collections = cmap.New()

type relationRepository struct {
	dao *Dao[*model.Relation]
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

func (r *relationRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *relationRepository) Update(ctx context.Context, tenantId string, v *model.Relation) error {
	return r.dao.Update(ctx, v)
}

func (r *relationRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *relationRepository) FindPaging(ctx context.Context, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.Relation], bool, error) {
	return r.dao.findPaging(ctx, query).Result()
}
