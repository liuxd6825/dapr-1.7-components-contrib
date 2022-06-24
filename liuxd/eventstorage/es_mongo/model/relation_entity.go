package model

import (
	"errors"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
)

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

func (r *RelationEntity) AddItem(idName, idValue string) {
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
