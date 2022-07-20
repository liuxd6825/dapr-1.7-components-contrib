package repository

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateRepository struct {
	BaseRepository
}

func NewAggregateRepository(client *mongo.Client, collection *mongo.Collection) *AggregateRepository {
	return &AggregateRepository{
		BaseRepository{
			client:     client,
			collection: collection,
		},
	}
}

func (r *AggregateRepository) FindById(ctx context.Context, tenantId string, aggregateId string) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	return r.findOne(ctx, filter)
}

func (r *AggregateRepository) Insert(ctx context.Context, aggregate *model.AggregateEntity) error {
	_, err := r.collection.InsertOne(ctx, aggregate)
	if err != nil {
		return err
	}
	return nil
}

//
// Delete
// @Description: 设置聚合为删除状态,并更新SequenceNumber
// @receiver r
// @param ctx
// @param tenantId
// @param aggregateId
// @return *model.AggregateEntity
// @return error
//
func (r *AggregateRepository) Delete(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	update := bson.M{
		"$set": bson.M{"deleted": true},
	}
	agg, err := r.findOneAndUpdate(ctx, aggregateId, filter, update)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

func (r *AggregateRepository) DeleteAndNextSequenceNumber(ctx context.Context, tenantId, aggregateId string) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	update := bson.M{
		"$set": bson.M{"deleted": true},
		"$inc": bson.M{SequenceNumberField: 1},
	}
	agg, err := r.findOneAndUpdate(ctx, aggregateId, filter, update)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

//
// NextSequenceNumber
// @Description: 获取新的消息序列号
// @receiver r
// @param ctx 上下文
// @param tenantId 租户ID
// @param aggregateId 聚合根Id
// @param count 新序列号的数量，单条消息时值为下1，多条消息时值为信息条数。
// @return *model.AggregateEntity 聚合对象
// @return error
//
func (r *AggregateRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, count uint64) (*model.AggregateEntity, error) {
	idValue, err := model.ObjectIDFromHex(aggregateId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		TenantIdField: tenantId,
		IdField:       idValue,
	}
	update := bson.M{
		"$inc": bson.M{SequenceNumberField: count},
	}
	agg, err := r.findOneAndUpdate(ctx, aggregateId, filter, update)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

func (r *AggregateRepository) findOneAndUpdate(ctx context.Context, aggregateId string, filter, update bson.M) (*model.AggregateEntity, error) {
	result := r.collection.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("aggregate idValue %s does not exist", aggregateId)
	} else if err != nil {
		return nil, err
	}
	var aggregate model.AggregateEntity
	if err := result.Decode(&aggregate); err != nil {
		return nil, err
	}
	return &aggregate, err
}

func (r *AggregateRepository) findOne(ctx context.Context, filter interface{}) (*model.AggregateEntity, error) {
	var result model.AggregateEntity
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
