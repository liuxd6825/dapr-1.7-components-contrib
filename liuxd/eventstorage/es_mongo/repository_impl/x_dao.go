package repository_impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/common/rsql"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

const (
	IdField             = "_id"
	TenantIdField       = "tenant_id"
	AggregateIdField    = "aggregate_id"
	AggregateTypeField  = "aggregate_type"
	EventIdField        = "event_id"
	SequenceNumberField = "sequence_number"
	PublishStatusField  = "publish_status"
)

type Options struct {
	DbId    string
	MaxTime *time.Duration
	Sort    interface{}
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetDbId(v string) *Options {
	o.DbId = v
	return o
}

func (o *Options) GetDbId() string {
	return o.DbId
}

func (o *Options) SetSort(d bson.D) *Options {
	o.Sort = d
	return o
}

func (o *Options) GetSort() interface{} {
	return o.Sort
}

func (o *Options) Merge(opts ...*Options) *Options {
	for _, item := range opts {
		if item == nil {
			continue
		}
		if len(item.DbId) != 0 {
			o.DbId = item.DbId
		}
		if item.MaxTime != nil {
			o.MaxTime = item.MaxTime
		}
	}
	return o
}

type dao[T model.Entity] struct {
	mongodb  *db.MongoDbConfig
	collName string
}

func NewDao[T model.Entity](mongodb *db.MongoDbConfig, collName string) *dao[T] {
	dao := dao[T]{
		mongodb:  mongodb,
		collName: collName,
	}
	return &dao
}

func (d *dao[T]) Save(ctx context.Context, v T, opts ...*Options) error {
	opt := NewOptions().SetDbId(v.GetTenantId()).Merge(opts...)

	filterMap := make(map[string]interface{})
	filterMap[IdField] = v.GetId()
	filter := d.NewFilter(v.GetTenantId(), filterMap)
	setData := bson.M{"$set": v}

	_, err := d.getCollection(opt.DbId).UpdateOne(ctx, filter, setData, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (d *dao[T]) DeleteById(ctx context.Context, tenantId, id string, opts ...*Options) error {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)

	filterMap := make(map[string]interface{})
	filterMap[IdField] = id
	filter := d.NewFilter(tenantId, filterMap)

	_, err := d.getCollection(opt.DbId).DeleteOne(ctx, filter, options.Delete())
	return err
}

func (d *dao[T]) DeleteById2(ctx context.Context, tenantId, id string, opts ...*Options) error {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	filterMap := make(map[string]interface{})
	filterMap[IdField] = id
	filter := d.NewFilter(tenantId, filterMap)
	_, err := d.getCollection(opt.DbId).DeleteOne(ctx, filter, options.Delete())
	return err
}

func (d *dao[T]) Delete(ctx context.Context, tenantId string, v T, opts ...*Options) error {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	return d.DeleteById(ctx, opt.DbId, v.GetId())
}

func (d *dao[T]) Insert(ctx context.Context, v T, opts ...*Options) error {
	opt := NewOptions().SetDbId(v.GetTenantId()).Merge(opts...)
	_, err := d.getCollection(opt.DbId).InsertOne(ctx, v)
	return err
}

func (d *dao[T]) InsertMany(ctx context.Context, tenantId string, vList []T, opts ...*Options) error {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	var docs []interface{}
	for _, rel := range vList {
		docs = append(docs, rel)
	}
	_, err := d.getCollection(opt.DbId).InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (d *dao[T]) Update(ctx context.Context, v T, opts ...*Options) error {
	opt := NewOptions().SetDbId(v.GetTenantId()).Merge(opts...)
	_, err := d.getCollection(opt.DbId).UpdateByID(ctx, v.GetId(), v)
	if err != nil {
		return err
	}
	return nil
}

func (d *dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...*Options) (T, bool, error) {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	filter := bson.D{
		{IdField, id},
		{TenantIdField, tenantId},
	}
	return d.findOne(ctx, opt.DbId, filter)
}

func (d *dao[T]) FindPaging(ctx context.Context, query dto.FindPagingQuery, opts ...*Options) *dto.FindPagingResult[T] {
	return d.findPaging(ctx, query, opts...)
}

func (d *dao[T]) findOneAndUpdate(ctx context.Context, tenantId string, filterData, updateData bson.M) (T, bool, error) {
	var null T
	result := d.getCollection(tenantId).FindOneAndUpdate(ctx, filterData, updateData)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return null, false, nil
	} else if err != nil {
		return null, false, err
	}

	var v T
	if err := result.Decode(&v); err != nil {
		return null, false, err
	}
	return v, true, err
}

func (d *dao[T]) findOne(ctx context.Context, tenantId string, filter interface{}, opts ...*Options) (T, bool, error) {
	var res T
	var null T
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	err := d.getCollection(opt.GetDbId()).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return null, false, nil
		}
		return null, false, err
	}
	return res, true, nil
}

func (d *dao[T]) deleteByFilter(ctx context.Context, tenantId string, filter interface{}, opts ...*Options) error {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	_, err := d.getCollection(opt.GetDbId()).DeleteMany(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

func (d *dao[T]) findList(ctx context.Context, tenantId string, filter bson.M, limit *int64, opts ...*Options) ([]T, bool, error) {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	findOpts := newFindOptions(opt)
	findOpts.Limit = limit
	cursor, err := d.getCollection(opt.GetDbId()).Find(ctx, filter, findOpts)
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	var list []T
	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, false, err
	}
	return list, len(list) > 0, nil
}

func (d *dao[T]) findPaging(ctx context.Context, query dto.FindPagingQuery, opts ...*Options) *dto.FindPagingResult[T] {
	return d.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*dto.FindPagingResult[T], bool, error) {
		var data []T
		opt := NewOptions().SetDbId(query.GetTenantId()).Merge(opts...)
		findOptions := &options.FindOptions{
			MaxTime: opt.MaxTime,
		}
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(int64(query.GetPageSize()))
			findOptions.SetSkip(int64(query.GetPageSize() * query.GetPageNum()))
		}
		if len(query.GetSort()) > 0 {
			sort, err := d.getSort(query.GetSort())
			if err != nil {
				return nil, false, err
			}
			findOptions.SetSort(sort)
		}

		coll := d.getCollection(opt.DbId)
		cursor, err := coll.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, &data)

		var totalRows int64 = -1
		if query.GetIsTotalRows() {
			totalRows, err = coll.CountDocuments(ctx, filter)
		}

		findData := dto.NewFindPagingResult[T](data, uint64(totalRows), query, err)
		return findData, true, err
	})

}

func (d *dao[T]) getCollection(dbId string) *mongo.Collection {
	collectionName := fmt.Sprintf("%v_%v", dbId, d.collName)
	value, ok := collections.Get(collectionName)
	if !ok {
		value = d.mongodb.NewCollection(collectionName)
		collections.Set(collectionName, value)
	}
	coll, _ := value.(*mongo.Collection)
	return coll
}

func (d *dao[T]) NewFilter(tenantId string, filterMap map[string]interface{}) bson.D {
	filter := bson.D{
		{TenantIdField, tenantId},
	}
	if filterMap != nil {
		for fieldName, fieldValue := range filterMap {
			if fieldName != IdField {
				fieldName = utils.AsMongoName(fieldName)
			}
			e := bson.E{
				Key:   fieldName,
				Value: fieldValue,
			}
			filter = append(filter, e)
		}
	}
	return filter
}

func (d *dao[T]) newBsonM(tenantId string, data map[string]interface{}) bson.M {
	return bson.M(data)
}

func (d *dao[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*dto.FindPagingResult[T], bool, error)) *dto.FindPagingResult[T] {
	p := NewMongoProcess()
	if err := rsql.ParseProcess(filter, p); err != nil {
		return dto.NewFindPagingResultWithError[T](err)
	}
	filterData := p.GetFilter(tenantId)
	data, _, err := fun(filterData)
	if err != nil {
		if IsErrorMongoNoDocuments(err) {
			err = nil
		}
	}
	return data
}

func (d *dao[T]) getSort(sort string) (map[string]interface{}, error) {
	if len(sort) == 0 {
		return nil, nil
	}
	//name:desc,id:asc
	res := map[string]interface{}{}
	list := strings.Split(sort, ",")
	for _, s := range list {
		sortItem := strings.Split(s, ":")
		name := sortItem[0]
		name = strings.Trim(name, " ")
		if name == "id" {
			name = IdField
		}
		order := "asc"
		if len(sortItem) > 1 {
			order = sortItem[1]
			order = strings.ToLower(order)
			order = strings.Trim(order, " ")
		}

		// 其中 1 为升序排列，而-1是用于降序排列.
		orderVal := 1
		var oerr error
		switch order {
		case "asc":
			orderVal = 1
		case "desc":
			orderVal = -1
		default:
			oerr = errors.New("order " + order + " is error")
		}
		if oerr != nil {
			return nil, oerr
		}
		res[name] = orderVal
	}
	return res, nil
}

func IsErrorMongoNoDocuments(err error) bool {
	if err == mongo.ErrNoDocuments {
		return true
	}
	return false
}

func newFindOptions(opt *Options) *options.FindOptions {
	findOptions := &options.FindOptions{}
	findOptions.MaxTime = opt.MaxTime
	findOptions.Sort = opt.Sort
	return findOptions
}
