package repository

import (
	"context"
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/rsql"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/other"
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

type BaseRepository[T any] struct {
	mongodb       *other.MongoDB
	collection    *mongo.Collection
	NewEntityList func() interface{}
}

func (r *BaseRepository[T]) FindPaging(ctx context.Context, collection *mongo.Collection, query eventstorage.FindPagingQuery, opts ...*other.FindOptions) *eventstorage.FindPagingResult[T] {
	return r.DoFilter(query.GetTenantId(), query.GetFilter(), func(filter map[string]interface{}) (*eventstorage.FindPagingResult[T], bool, error) {
		data := r.NewEntityList()
		findOptions := getFindOptions(opts...)
		if query.GetPageSize() > 0 {
			findOptions.SetLimit(int64(query.GetPageSize()))
			findOptions.SetSkip(int64(query.GetPageSize() * query.GetPageNum()))
		}
		if len(query.GetSort()) > 0 {
			sort, err := r.getSort(query.GetSort())
			if err != nil {
				return nil, false, err
			}
			findOptions.SetSort(sort)
		}

		cursor, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, false, err
		}
		err = cursor.All(ctx, data)
		totalRows, err := collection.CountDocuments(ctx, filter)
		findData := eventstorage.NewFindPagingResult[T](data.(*[]T), uint64(totalRows), query, err)
		return findData, true, err
	})

}

func (r *BaseRepository[T]) NewFilter(tenantId string, filterMap map[string]interface{}) bson.D {
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

func (r *BaseRepository[T]) DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*eventstorage.FindPagingResult[T], bool, error)) *eventstorage.FindPagingResult[T] {
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

func (r *BaseRepository[T]) getSort(sort string) (map[string]interface{}, error) {
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

func getFindOptions(opts ...*other.FindOptions) *options.FindOptions {
	opt := MergeFindOptions(opts...)
	findOneOptions := &options.FindOptions{}
	findOneOptions.MaxTime = opt.MaxTime
	return findOneOptions
}

func MergeFindOptions(opts ...*other.FindOptions) *other.FindOptions {
	res := &other.FindOptions{}
	for _, o := range opts {
		if o.MaxTime != nil {
			res.MaxTime = o.MaxTime
		}
	}
	return res
}
