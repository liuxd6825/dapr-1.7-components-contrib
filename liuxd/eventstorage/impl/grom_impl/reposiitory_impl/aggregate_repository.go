package repository_impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type aggregateRepository struct {
	dao *dao[*model.Aggregate]
}

func NewAggregateRepository(db *gorm.DB) repository.AggregateRepository {
	_ = db.AutoMigrate(&model.Aggregate{})
	return &aggregateRepository{
		dao: NewDao[*model.Aggregate](db,
			func() *model.Aggregate { return &model.Aggregate{} },
			func() []*model.Aggregate { return []*model.Aggregate{} },
		),
	}
}

func (r *aggregateRepository) Create(ctx context.Context, v *model.Aggregate) error {
	return r.dao.Insert(ctx, v)
}

func (r *aggregateRepository) DeleteById(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *aggregateRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	where := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v"`, tenantId, aggregateId)
	return r.dao.deleteByFilter(ctx, tenantId, where)
}

func (r *aggregateRepository) Update(ctx context.Context, v *model.Aggregate) error {
	return r.dao.Update(ctx, v)
}

func (r *aggregateRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Aggregate, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *aggregateRepository) UpdateIsDelete(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	update := map[string]interface{}{
		"deleted": true,
	}
	if agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, aggregateId, update); err != nil {
		return nil, false, err
	} else {
		return agg, ok, nil
	}
}

//
// SetIsDelete
// @Description: 设置聚合为删除状态,并更新SequenceNumber
// @receiver r
// @param ctx
// @param tenantId
// @param aggregateId
// @return *model.Aggregate
// @return error
//
func (r *aggregateRepository) SetIsDelete(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	update := map[string]interface{}{
		"$set": bson.M{"deleted": true},
	}
	agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, aggregateId, update)
	return agg, ok, err
}

func (r *aggregateRepository) DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	sql := fmt.Sprintf("update aggregates set deleted=true, sequence_number=sequence_number+1 where tanent_id=%v and id=%v", tenantId, aggregateId)
	if err := r.dao.execSql(ctx, sql); err != nil {
		return nil, false, err
	}
	agg, ok, err := r.dao.FindById(ctx, tenantId, aggregateId)
	return agg, ok, err
}

//
// NextSequenceNumber
// @Description: 获取新的消息序列号
// @receiver r
// @param ctx 上下文
// @param tenantId 租户ID
// @param aggregateId 聚合根Id
// @param count 新序列号的数量，单条消息时值为下1，多条消息时值为信息条数。
// @return *model.Aggregate 聚合对象
// @return error
//
func (r *aggregateRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, count uint64) (*model.Aggregate, bool, uint64, error) {
	sql := fmt.Sprintf(`update aggregates set sequence_number=sequence_number+%v where tenant_id="%v" and id="%v"`, count, tenantId, aggregateId)
	err := r.dao.execSql(ctx, sql)
	if err != nil {
		return nil, false, 0, err
	}
	agg, ok, err := r.dao.FindById(ctx, tenantId, aggregateId)
	if !ok {
		return agg, false, 0, nil
	}
	return agg, ok, agg.SequenceNumber + 1, nil
}
