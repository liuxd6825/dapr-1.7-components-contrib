package model

type Aggregate struct {
	Id             string `bson:"_id" json:"id" gorm:"primaryKey"`
	TenantId       string `bson:"tenant_id" json:"tenant_id"`
	AggregateId    string `bson:"aggregate_id" json:"aggregate_id"`
	AggregateType  string `bson:"aggregate_type" json:"aggregate_type"`
	SequenceNumber uint64 `bson:"sequence_number" json:"sequence_number"`
	Deleted        bool   `bson:"deleted" json:"deleted"`
}

func (a *Aggregate) GetId() string {
	return a.Id
}

func (a *Aggregate) SetId(v string) {
	a.Id = v
}

func (a *Aggregate) GetTenantId() string {
	return a.TenantId
}

type AggregateBuilder struct {
}

func NewAggregateBuilder() EntityBuilder {
	return &AggregateBuilder{}
}

func (a AggregateBuilder) NewEntity() interface{} {
	return &Aggregate{}
}

func (a AggregateBuilder) NewEntities() interface{} {
	return []*Aggregate{}
}
