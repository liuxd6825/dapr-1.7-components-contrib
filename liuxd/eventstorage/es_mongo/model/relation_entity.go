package model

import "github.com/liuxd6825/components-contrib/liuxd/common/utils"

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
	for key, value := range items {
		res.AddItem(key, value)
	}
	return res
}

func (r *RelationEntity) AddItem(idName, idValue string) {
	name := utils.AsMongoName(idName)
	r.Items[name] = idValue
}
