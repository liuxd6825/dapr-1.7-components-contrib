package repository_impl

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
)

type aggregateRepository struct {
	dao *dao[*model.Aggregate]
}

func NewAggregateRepository(mongodb *db.MongoDbConfig, collName string) repository.AggregateRepository {
	return &aggregateRepository{
		dao: NewDao[*model.Aggregate](mongodb, collName),
	}
}

func (r *aggregateRepository) Create(ctx context.Context, v *model.Aggregate) error {
	return r.dao.Insert(ctx, v)
}

func (r *aggregateRepository) Delete(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *aggregateRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	filter := bson.M{
		TenantIdField:    tenantId,
		AggregateIdField: aggregateId,
	}
	return r.dao.deleteByFilter(ctx, tenantId, filter)
}

func (r *aggregateRepository) Update(ctx context.Context, v *model.Aggregate) error {
	return r.dao.Update(ctx, v)
}

func (r *aggregateRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Aggregate, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *aggregateRepository) UpdateIsDelete(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	filter := map[string]interface{}{
		TenantIdField: tenantId,
		IdField:       aggregateId,
	}
	update := map[string]interface{}{
		"$set": bson.M{"deleted": true},
	}
	agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, filter, update)
	return agg, ok, err
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
	filter := map[string]interface{}{
		TenantIdField: tenantId,
		IdField:       aggregateId,
	}
	update := map[string]interface{}{
		"$set": bson.M{"deleted": true},
	}
	agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, filter, update)
	return agg, ok, err
}

func (r *aggregateRepository) DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.Aggregate, bool, error) {
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       aggregateId,
	}
	update := bson.M{
		"$set": bson.M{"deleted": true},
		"$inc": bson.M{SequenceNumberField: 1},
	}
	agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, filter, update)
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
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       aggregateId,
	}
	update := bson.M{
		"$inc": bson.M{SequenceNumberField: count},
	}
	agg, ok, err := r.dao.findOneAndUpdate(ctx, tenantId, filter, update)
	if err != nil {
		return nil, ok, 0, err
	}
	if !ok {
		return agg, false, 0, nil
	}
	return agg, ok, agg.SequenceNumber + 1, nil
}
