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
	Id            string `bson:"_id" json:"id"  gorm:"primaryKey"`
	TenantId      string `bson:"tenant_id" json:"tenant_id"`
	TableName     string `bson:"table_name" json:"table_name"`
	AggregateId   string `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType string `bson:"aggregate_type" json:"aggregate_type"`
	EventId       string `bson:"event_id" json:"event_id"`
	EventType     string `bson:"event_type" json:"event_type"`
	IsDeleted     bool   `bson:"is_deleted" json:"is_deleted"`
	RelName       string `bson:"rel_name" json:"rel_name"`
	RelValue      string `bson:"rel_value" json:"rel_value"`
}

func NewRelations(tenantId, eventId, eventType, aggregateId, aggregateType string, items map[string]string) []*Relation {
	tableName := utils.AsMongoName(aggregateType)
	var relations []*Relation
	for relName, relValue := range items {
		if len(relValue) == 0 {
			continue
		}
		rel := &Relation{
			Id:            NewObjectID(),
			TenantId:      tenantId,
			AggregateId:   aggregateId,
			AggregateType: aggregateType,
			EventId:       eventId,
			EventType:     eventType,
			TableName:     tableName,
			IsDeleted:     false,
			RelName:       relName,
			RelValue:      relValue,
		}
		relations = append(relations, rel)
	}
	return relations
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
	if len(r.AggregateType) == 0 {
		return errors.New("Relation.AggregateType cannot be empty")
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
