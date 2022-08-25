package repository_impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/common/rsql"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
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
)

type Options struct {
	DbId    string
	MaxTime *time.Duration
	Sort    *string
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

func (o *Options) SetSort(v *string) *Options {
	o.Sort = v
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
		if item.Sort != nil {
			o.Sort = item.Sort
		}
	}
	return o
}

type dao[T model.Entity] struct {
	db          *gorm.DB
	newFunc     func() T
	newListFunc func() []T
}

func NewDao[T model.Entity](db *gorm.DB, newFunc func() T, newListFunc func() []T) *dao[T] {
	dao := dao[T]{
		db:          db,
		newFunc:     newFunc,
		newListFunc: newListFunc,
	}
	return &dao
}

func (d *dao[T]) NewEntity() T {
	return d.newFunc()
}

func (d *dao[T]) NewEntities() []T {
	return d.newListFunc()
}

func (d *dao[T]) Save(ctx context.Context, v T, opts ...*Options) error {
	return d.getDB(ctx).Model(&v).Updates(v).Error
}

func (d *dao[T]) DeleteById(ctx context.Context, tenantId, id string, opts ...*Options) error {
	v := d.NewEntity()
	return d.getDB(ctx).Where("tenant_id=? and id=?", tenantId, id).Delete(v).Error
}

func (d *dao[T]) Delete(ctx context.Context, tenantId string, v T, opts ...*Options) error {
	return d.getDB(ctx).Model(v).Delete(v).Error
}

func (d *dao[T]) Insert(ctx context.Context, v T, opts ...*Options) error {
	return d.getDB(ctx).Model(v).Create(v).Error
}

func (d *dao[T]) InsertMany(ctx context.Context, tenantId string, vList []T, opts ...*Options) error {
	v := d.NewEntity()
	return d.getDB(ctx).Model(v).CreateInBatches(vList, len(vList)).Error
}

func (d *dao[T]) Update(ctx context.Context, v T, opts ...*Options) error {
	return d.getDB(ctx).Where("tenant_id=? and id=?", v.GetTenantId(), v.GetId()).Updates(v).Error
}

func (d *dao[T]) UpdateByMap(ctx context.Context, tenantId, id string, data map[string]interface{}, opts ...*Options) error {
	var v T = d.NewEntity()
	return d.getDB(ctx).Model(v).Where("tenant_id=? and id=?", tenantId, id).Updates(data).Error
}

func (d *dao[T]) FindById(ctx context.Context, tenantId string, id string, opts ...*Options) (T, bool, error) {
	var null T
	entity := d.NewEntity()
	err := d.getDB(ctx).Where("tenant_id=? and id=?", tenantId, id).First(entity).Error
	if IsErrRecordNotFound(err) {
		return null, false, nil
	} else if err != nil {
		return null, false, err
	}
	return entity, true, nil
}

func (d *dao[T]) FindPaging(ctx context.Context, query dto.FindPagingQuery, opts ...*Options) *dto.FindPagingResult[T] {
	return d.findPaging(ctx, query, opts...)
}

func (d *dao[T]) execSql(ctx context.Context, sql string) error {
	err := d.getDB(ctx).Exec(sql).Error
	if IsErrRecordNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (d *dao[T]) findOneAndUpdate(ctx context.Context, tenantId string, id string, updateSet map[string]interface{}) (T, bool, error) {
	if err := d.UpdateByMap(ctx, tenantId, id, updateSet); err != nil {
		var null T
		return null, false, err
	}
	agg, ok, err := d.FindById(ctx, tenantId, id)
	return agg, ok, err
}

func (d *dao[T]) findOne(ctx context.Context, tenantId string, where string, opts ...*Options) (T, bool, error) {
	res := d.NewEntity()
	var null T
	err := d.getDB(ctx).Model(res).Where(where).First(res).Error
	if IsErrRecordNotFound(err) {
		return null, false, nil
	} else if err != nil {
		return null, false, err
	}
	return res, true, nil
}

func (d *dao[T]) deleteByFilter(ctx context.Context, tenantId string, where string, opts ...*Options) error {
	return d.getDB(ctx).Where(where).Delete(d.NewEntity()).Error
}

func (d *dao[T]) findList(ctx context.Context, tenantId string, where string, limit *int64, opts ...*Options) ([]T, bool, error) {
	opt := NewOptions().SetDbId(tenantId).Merge(opts...)
	findOpts := newFindOptions(opt)
	findOpts.Limit = limit
	list := d.NewEntities()
	tx := d.getDB(ctx)
	var err error
	if limit != nil {
		var l int = 0
		l = int(*limit)
		tx = tx.Limit(l)
	}
	if opt.Sort != nil {
		order := *opt.Sort
		tx = tx.Order(order)
	}
	w := "tenant_id=?"
	if len(where) > 0 {
		w = fmt.Sprintf("(%v) and tenant_id=?", where)
	}
	err = tx.Where(w, tenantId).Find(&list).Error
	if IsErrRecordNotFound(err) {
		return list, false, nil
	} else if err != nil {
		return list, false, err
	}
	return list, true, nil
}

func (d *dao[T]) findPaging(ctx context.Context, query dto.FindPagingQuery, opts ...*Options) *dto.FindPagingResult[T] {
	return d.DoFilter(query.GetTenantId(), query.GetFilter(), func(sqlWhere string) (*dto.FindPagingResult[T], bool, error) {
		var data []T = d.NewEntities()

		tx := d.getDB(ctx)

		if len(sqlWhere) > 0 {
			tx = tx.Where(sqlWhere)
		}

		if query.GetPageSize() > 0 {
			tx = tx.Limit(int(query.GetPageSize()))
		}

		if query.GetPageNum() > 0 {
			tx = tx.Offset(int(query.GetPageSize() * query.GetPageNum()))
		}

		if len(query.GetSort()) > 0 {
			tx = tx.Order(query.GetSort())
		}

		err := tx.Find(&data).Error
		if IsErrRecordNotFound(err) {
			return nil, false, nil
		}

		var totalRows int64 = -1
		if query.GetIsTotalRows() {
			tx := d.getDB(ctx)
			if len(sqlWhere) > 0 {
				tx = tx.Where(sqlWhere)
			}
			if err := tx.Count(&totalRows).Error; err != nil {
				return nil, false, err
			}
		}

		findData := dto.NewFindPagingResult[T](data, uint64(totalRows), query, err)
		return findData, true, err
	})

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

func (d *dao[T]) DoFilter(tenantId, filter string, fun func(sqlWhere string) (*dto.FindPagingResult[T], bool, error)) *dto.FindPagingResult[T] {
	p := rsql.NewSqlProcess()
	if err := rsql.ParseProcess(filter, p); err != nil {
		return dto.NewFindPagingResultWithError[T](err)
	}
	sqlWhere := p.GetFilter(tenantId)
	data, _, err := fun(sqlWhere.(string))
	if err != nil {
		if IsErrRecordNotFound(err) {
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

func (d *dao[T]) getDB(ctx context.Context) *gorm.DB {
	tx := db.GetTransaction(ctx)
	if tx != nil {
		return tx
	}
	return d.db
}

func IsErrRecordNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
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
