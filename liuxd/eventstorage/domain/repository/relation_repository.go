package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
)

type RelationRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Relation) error
	CreateMany(ctx context.Context, tenantId string, vList []*model.Relation) error
	DeleteById(ctx context.Context, tenantId string, id string) error
	DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error
	Update(ctx context.Context, tenantId string, v *model.Relation) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error)
	FindPaging(ctx context.Context, query dto.FindPagingQuery) (*dto.FindPagingResult[*model.Relation], bool, error)
}
