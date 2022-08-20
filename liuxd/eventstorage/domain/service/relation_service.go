package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type RelationService interface {
	Save(ctx context.Context, relation *model.Relation) error
	FindPaging(ctx context.Context, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.Relation], bool, error)
}

func NewRelationService(repos repository.RelationRepository) RelationService {
	return &relationService{repos: repos}
}

type relationService struct {
	repos repository.RelationRepository
}

func (r *relationService) Create(ctx context.Context, relation *model.Relation) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.repos.Create(ctx, relation.TenantId, relation)
}

func (r *relationService) Update(ctx context.Context, relation *model.Relation) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.repos.Update(ctx, relation.TenantId, relation)
}

func (r *relationService) Delete(ctx context.Context, tenantId string, id string) error {
	return r.repos.Delete(ctx, tenantId, id)
}

func (r *relationService) FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error) {
	return r.repos.FindById(ctx, tenantId, id)
}

func (r *relationService) Save(ctx context.Context, relation *model.Relation) error {
	if err := relation.Validate(); err != nil {
		return err
	}
	return r.repos.Update(ctx, relation.TenantId, relation)
}

func (r *relationService) FindPaging(ctx context.Context, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.Relation], bool, error) {
	res, ok, err := r.repos.FindPaging(ctx, query)
	return res, ok, err
}
