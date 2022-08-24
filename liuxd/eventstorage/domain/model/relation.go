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

type RelationItems map[string]string

type Relation struct {
	Id          string        `bson:"_id" json:"id"  gorm:"primaryKey"`
	TenantId    string        `bson:"tenant_id" json:"tenant_id"`
	TableName   string        `bson:"table_name" json:"table_name"`
	AggregateId string        `bson:"aggregate_id" json:"aggregate_id"`
	IsDeleted   bool          `bson:"is_deleted" json:"is_deleted"`
	Items       RelationItems `bson:",inline" json:"items"`
}

func NewRelationEntity(tenantId, relationId, aggregateId, aggregateType string, items map[string]string) *Relation {
	tableName := utils.AsMongoName(aggregateType)
	res := &Relation{
		Id:          relationId,
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

func (r Relation) AddItem(idName, idValue string) {
	name := utils.AsMongoName(idName)
	r.Items[name] = idValue
}

func (r *Relation) Validate() error {
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

func (r *Relation) GetId() string {
	return r.Id
}

func (r *Relation) SetId(v string) {
	r.Id = v
}

func (r *Relation) GetTenantId() string {
	return r.TenantId
}
