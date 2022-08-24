package repository_impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"gorm.io/gorm"
)

type snapshotRepository struct {
	dao *dao[*model.Snapshot]
}

func NewSnapshotRepository(db *gorm.DB) repository.SnapshotRepository {
	_ = db.AutoMigrate(&model.Snapshot{})
	return &snapshotRepository{
		dao: NewDao[*model.Snapshot](db,
			func() *model.Snapshot { return &model.Snapshot{} },
			func() []*model.Snapshot { return []*model.Snapshot{} },
		),
	}
}

func (r *snapshotRepository) Create(ctx context.Context, tenantId string, v *model.Snapshot) error {
	return r.dao.Insert(ctx, v)
}

func (r *snapshotRepository) DeleteById(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *snapshotRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	where := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v"`, tenantId, aggregateId)
	return r.dao.deleteByFilter(ctx, tenantId, where)
}

func (r *snapshotRepository) Update(ctx context.Context, tenantId string, v *model.Snapshot) error {
	return r.dao.Update(ctx, v)
}

func (r *snapshotRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Snapshot, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *snapshotRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) ([]*model.Snapshot, bool, error) {
	filter := fmt.Sprintf(`aggregate_id="%v"`, aggregateId)
	return r.dao.findList(ctx, tenantId, filter, nil)
}

func (r *snapshotRepository) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.Snapshot, bool, error) {
	filter := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v" and aggregate_type="%v"`, tenantId, aggregateId, aggregateType)
	sort := fmt.Sprintf("sequence_number asc")
	options := NewOptions().SetSort(&sort)
	return r.dao.findOne(ctx, tenantId, filter, options)
}
