package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
)

type RelationService interface {
	CreateMany(ctx context.Context, tenantId string, relation []*model.Relation) error
	DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error
	FindPaging(ctx context.Context, query dto.FindPagingQuery) (*dto.FindPagingResult[*model.Relation], bool, error)
}

func NewRelationService(repos repository.RelationRepository) RelationService {
	return &relationService{repos: repos}
}

type relationService struct {
	repos repository.RelationRepository
}

func (r *relationService) Create(ctx context.Context, relation *model.Relation) error {
	if relation == nil {
		return nil
	}
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.repos.Create(ctx, relation.TenantId, relation)
}

func (r *relationService) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	return r.repos.DeleteByAggregateId(ctx, tenantId, aggregateId)
}

func (r *relationService) Update(ctx context.Context, relation *model.Relation) error {
	if relation == nil {
		return nil
	}
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.repos.Update(ctx, relation.TenantId, relation)
}

func (r *relationService) Delete(ctx context.Context, tenantId string, id string) error {
	return r.repos.DeleteById(ctx, tenantId, id)
}

func (r *relationService) FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error) {
	return r.repos.FindById(ctx, tenantId, id)
}

func (r *relationService) CreateMany(ctx context.Context, tenantId string, relations []*model.Relation) error {
	if len(relations) == 0 {
		return nil
	}
	return r.repos.CreateMany(ctx, tenantId, relations)
}

func (r *relationService) FindPaging(ctx context.Context, query dto.FindPagingQuery) (*dto.FindPagingResult[*model.Relation], bool, error) {
	res, ok, err := r.repos.FindPaging(ctx, query)
	return res, ok, err
}
