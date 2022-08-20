package repository_impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/common/rsql"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
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

type Dao[T model.Entity] struct {
	mongodb       *db.MongoDbConfig
	collName      string
	NewEntityList func() interface{}
}

func NewDao[T model.Entity](mongodb *db.MongoDbConfig, collName string) *Dao[T] {
	dao := Dao[T]{
		mongodb:  mongodb,
		collName: collName,
	}
	return &dao
}

func (d *Dao[T]) Save(ctx context.Context, v T) error {
	filterMap := make(map[string]interface{})
	entity := d.asEntity(v)
	filterMap[IdField] = entity.GetId()
	filter := d.NewFilter(entity.GetTenantId(), filterMap)
	setData := bson.M{"$set": v}
	_, err := d.getCollection(entity.GetTenantId()).UpdateOne(ctx, filter, setData, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) asEntity(v interface{}) model.Entity {
	e, _ := v.(model.Entity)
	return e
}

func (d *Dao[T]) DeleteById(ctx context.Context, tenantId, id string) error {
	filterMap := make(map[string]interface{})
	filterMap[IdField] = id
	filter := d.NewFilter(tenantId, filterMap)
	_, err := d.getCollection(tenantId).DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) DeleteById2(ctx context.Context, tenantId, id string) error {
	filterMap := make(map[string]interface{})
	filterMap[IdField] = id
	filter := d.NewFilter(tenantId, filterMap)
	_, err := d.getCollection(tenantId).DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) Delete(ctx context.Context, v T) error {
	entity := d.asEntity(v)
	return d.DeleteById(ctx, entity.GetTenantId(), entity.GetId())
}

func (d *Dao[T]) Insert(ctx context.Context, v T) error {
	entity := d.asEntity(v)
	_, err := d.getCollection(entity.GetTenantId()).InsertOne(ctx, v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) InsertMany(ctx context.Context, tenantId string, vList []T) error {
	var docs []interface{}
	for _, rel := range vList {
		docs = append(docs, rel)
	}
	_, err := d.getCollection(tenantId).InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) Update(ctx context.Context, v T) error {
	entity := d.asEntity(v)
	_, err := d.getCollection(entity.GetTenantId()).UpdateByID(ctx, entity.GetId(), v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao[T]) FindById(ctx context.Context, tenantId string, id string) (T, bool, error) {
	filter := bson.D{
		{IdField, id},
		{TenantIdField, tenantId},
	}
	return d.findOne(ctx, tenantId, filter)
}

func (d *Dao[T]) FindPaging(ctx context.Context, query eventstorage.FindPagingQuery, opts ...*options.FindOptions) *eventstorage.FindPagingResult[T] {
	return d.findPaging(ctx, query, opts...)
}

func (d *Dao[T]) findOneAndUpdate(ctx context.Context, tenantId string, filterData, updateData bson.M) (T, bool, error) {
	var null T
	result := d.getCollection(tenantId).FindOneAndUpdate(ctx, filterData, updateData)
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return null, false, nil
	} else if err != nil {
		return null, false, err
	}
	var v T
	if err := result.Decode(v); err != nil {
		return null, false, err
	}
	return v, true, err
}

func (d *Dao[T]) findOne(ctx context.Context, tenantId string, filter interface{}) (T, bool, error) {
	var res T
	var null T
	err := d.getCollection(tenantId).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return null, false, nil
		}
		return null, false, err
	}

	return res, true, nil
}

func (d *Dao[T]) findList(ctx context.Context, tenantId string, filter bson.M, findOptions ...*options.FindOptions) ([]T, bool, error) {
	cursor, err := d.getCollection(tenantId).Find(ctx, filter, findOptions...)
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

func (d *Dao[T]) findPaging(ctx context.Context, query eventstorage.FindPagingQuery, opts ...*options.FindOptions) *eventstorage.FindPagingResult[T] {
	return d.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*eventstorage.FindPagingResult[T], bool, error) {
		var data []T
		findOptions := getFindOptions(opts...)
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

		coll := d.getCollection(query.GetTenantId())
		cursor, err := coll.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}

		err = cursor.All(ctx, &data)
		totalRows, err := coll.CountDocuments(ctx, filter)
		findData := eventstorage.NewFindPagingResult[T](data, uint64(totalRows), query, err)
		return findData, true, err
	})

}

func (d *Dao[T]) getCollection(tenantId string) *mongo.Collection {
	collectionName := fmt.Sprintf("%v_%v", tenantId, d.collName)
	value, ok := collections.Get(collectionName)
	if !ok {
		value = d.mongodb.NewCollection(collectionName)
		collections.Set(collectionName, value)
	}
	coll, _ := value.(*mongo.Collection)
	return coll
}

func (d *Dao[T]) NewFilter(tenantId string, filterMap map[string]interface{}) bson.D {
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

func (d *Dao[T]) newBsonM(tenantId string, data map[string]interface{}) bson.M {
	return bson.M(data)
}

func (d *Dao[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*eventstorage.FindPagingResult[T], bool, error)) *eventstorage.FindPagingResult[T] {
	p := NewMongoProcess()
	if err := rsql.ParseProcess(filter, p); err != nil {
		return eventstorage.NewFindPagingResultWithError[T](err)
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

func (d *Dao[T]) getSort(sort string) (map[string]interface{}, error) {
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

func getFindOptions(opts ...*options.FindOptions) *options.FindOptions {
	opt := MergeFindOptions(opts...)
	findOneOptions := &options.FindOptions{}
	findOneOptions.MaxTime = opt.MaxTime
	return findOneOptions
}

func MergeFindOptions(opts ...*options.FindOptions) *options.FindOptions {
	res := &options.FindOptions{}
	for _, o := range opts {
		if o.MaxTime != nil {
			res.MaxTime = o.MaxTime
		}
	}
	return res
}
