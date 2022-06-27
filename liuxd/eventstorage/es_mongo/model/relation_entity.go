package model

import (
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
)

const (
	RelationIdField     = "_id"
	RelationTenantId    = "tenant_id"
	RelationTableName   = "table_name"
	RelationAggregateId = "aggregate_id"
	RelationIsDeleted   = "is_deleted"
)

/*
type RelationEntity map[string]interface{}

func (r RelationEntity) SetId(v string) {
	r[RelationIdField] = v
}

func (r RelationEntity) GetId() string {
	return r.GetString(RelationIdField)
}

func (r RelationEntity) SetTenantId(v string) {
	r[RelationTenantId] = v
}

func (r RelationEntity) GetTenantId() string {
	return r.GetString(RelationTenantId)
}

func (r RelationEntity) SetTableName(v string) {
	r[RelationTableName] = v
}

func (r RelationEntity) GetTableName() string {
	return r.GetString(RelationTableName)
}

func (r RelationEntity) SetAggregateId(v string) {
	r[RelationAggregateId] = v
}

func (r RelationEntity) GetTableAggregateId() string {
	return r.GetString(RelationAggregateId)
}

func (r RelationEntity) SetIsDeleted(v bool) {
	r[RelationIsDeleted] = v
}

func (r RelationEntity) GetIsDeleted() bool {
	return r.GetBool(RelationIsDeleted)
}

func (r RelationEntity) AddItem(idName, idValue string) {
	name := utils.AsMongoName(idName)
	r[name] = idValue
}

func (r RelationEntity) GetBool(key string) bool {
	v, ok := r[key]
	if ok {
		return v.(bool)
	}
	return false
}

func (r RelationEntity) GetString(key string) string {
	v, ok := r[key]
	if ok {
		return v.(string)
	}
	return ""
}
*/

type RelationItems map[string]string

type RelationEntity struct {
	Id          string        `bson:"_id"`
	TenantId    string        `bson:"tenant_id"`
	TableName   string        `bson:"table_name"`
	AggregateId string        `bson:"aggregate_id"`
	IsDeleted   bool          `bson:"is_deleted"`
	Items       RelationItems `bson:",inline"`
}

func NewRelationEntity(tenantId, aggregateId, aggregateType string, items map[string]string) *RelationEntity {
	tableName := utils.AsMongoName(aggregateType)
	/*
		res := &RelationEntity{}
		res.SetTableName(tableName)
		res.SetAggregateId(aggregateId)
		res.SetAggregateId(aggregateId)
		res.SetTenantId(tenantId)
		res.SetIsDeleted(false)
	*/
	res := &RelationEntity{
		Id:          aggregateId,
		TenantId:    tenantId,
		AggregateId: aggregateId,
		TableName:   tableName,
		IsDeleted:   false,
		Items:       map[string]string{},
	}
	// 添加关系。注意：如果关系值是空，则不添加。
	if items != nil && len(items) > 0 {
		for key, value := range items {
			if len(value) > 0 {
				res.AddItem(key, value)
			}
		}
	}
	return res
}

func (r RelationEntity) AddItem(idName, idValue string) {
	name := utils.AsMongoName(idName)
	r.Items[name] = idValue
}

func (r *RelationEntity) Validate() error {
	if r == nil {
		return errors.New("relation is nil")
	}
	if len(r.TableName) == 0 {
		return errors.New("Relation.TableName cannot be empty")
	}
	if len(r.TenantId) == 0 {
		return errors.New("Relation.TenantId cannot be empty")
	}
	if len(r.AggregateId) == 0 {
		return errors.New("Relation.AggregateId cannot be empty")
	}
	if len(r.Id) == 0 {
		return errors.New("Relation.Id cannot be empty")
	}
	return nil
}
