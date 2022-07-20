package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"

	//"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/repository"
)

type RelationService interface {
	Save(ctx context.Context, relation *model.RelationEntity) error
	FindPaging(ctx context.Context, tableName string, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.RelationEntity], bool, error)
}

func NewRelationService(db *db.MongoDB) RelationService {
	res := &relationService{}
	res.resp = repository.NewRelationRepository(db)
	return res
}

type relationService struct {
	resp *repository.RelationRepository
}

func (r *relationService) Save(ctx context.Context, relation *model.RelationEntity) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.resp.Save(ctx, relation)
}

func (r *relationService) Create(ctx context.Context, relation *model.RelationEntity) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.resp.InsertOne(ctx, relation)
}

func (r *relationService) Update(ctx context.Context, relation *model.RelationEntity) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.resp.UpdateOne(ctx, relation)
}

func (r *relationService) FindPaging(ctx context.Context, tableName string, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.RelationEntity], bool, error) {
	res, ok, err := r.resp.FindPaging(ctx, tableName, query).Result()
	return res, ok, err
}
