package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
)

type RelationRepository interface {
	Create(ctx context.Context, tenantId string, v *model.Relation) error
	Delete(ctx context.Context, tenantId string, id string) error
	Update(ctx context.Context, tenantId string, v *model.Relation) error
	FindById(ctx context.Context, tenantId string, id string) (*model.Relation, bool, error)
	FindPaging(ctx context.Context, query eventstorage.FindPagingQuery) (*eventstorage.FindPagingResult[*model.Relation], bool, error)
}
