package repository_impl

import (
	"context"
	"fmt"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/dto"
	"gorm.io/gorm"
)

type eventRepository struct {
	dao *dao[*model.Event]
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	_ = db.AutoMigrate(&model.Event{})
	return &eventRepository{
		dao: NewDao[*model.Event](db,
			func() *model.Event { return &model.Event{} },
			func() []*model.Event { return []*model.Event{} },
		),
	}
}

func (r eventRepository) Create(ctx context.Context, tenantId string, v *model.Event) error {
	return r.dao.Insert(ctx, v)
}

func (r eventRepository) DeleteById(ctx context.Context, tenantId string, id string) error {
	return r.dao.DeleteById(ctx, tenantId, id)
}

func (r *eventRepository) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	where := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v"`, tenantId, aggregateId)
	return r.dao.deleteByFilter(ctx, tenantId, where)
}

func (r eventRepository) Update(ctx context.Context, tenantId string, v *model.Event) error {
	return r.dao.Update(ctx, v)
}

func (r eventRepository) FindById(ctx context.Context, tenantId string, id string) (*model.Event, bool, error) {
	return r.dao.FindById(ctx, tenantId, id)
}

func (r *eventRepository) FindPaging(ctx context.Context, query dto.FindPagingQuery) *dto.FindPagingResult[*model.Event] {
	return r.dao.FindPaging(ctx, query)
}

func (r *eventRepository) FindByEventId(ctx context.Context, tenantId string, eventId string) (*model.Event, bool, error) {
	filter := fmt.Sprintf(`tenant_id="%v" and event_id="%v"`, tenantId, eventId)
	return r.dao.findOne(ctx, tenantId, filter)
}

func (r *eventRepository) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string, aggregateType string) ([]*model.Event, bool, error) {
	filter := fmt.Sprintf(`tenant_id="%v" and aggregate_id="%v"`, tenantId, aggregateId)
	return r.dao.findList(ctx, tenantId, filter, nil)
}

//
// FindByGtSequenceNumber
// @Description: 查找大于SequenceNumber的事件
// @receiver r
// @param ctx
// @param tenantId
// @param aggregateId
// @param aggregateType
// @param sequenceNumber
// @return []*model.Event
// @return bool
// @return error
//
func (r *eventRepository) FindByGtSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string, sequenceNumber uint64) ([]*model.Event, bool, error) {
	filter := fmt.Sprintf(`aggregate_id="%v" and aggregate_type="%v"`, aggregateId, aggregateType)
	sort := fmt.Sprintf("%v asc", SequenceNumberField)
	findOptions := NewOptions().SetSort(&sort)
	return r.dao.findList(ctx, tenantId, filter, nil, findOptions)
}

/*
func (r *eventRepository) NextSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) uint64 {
	filter := bson.M{
		TenantIdField:      tenantId,
		AggregateIdField:   aggregateId,
		AggregateTypeField: aggregateType,
	}
	findOptions := options.FindOne().SetSort(bson.D{{SequenceNumberField, -1}})
	result := r.collection.FindOne(ctx, filter, findOptions)
	var event model.Event
	if err := result.Decode(&event); err == nil {
		return event.SequenceNumber + 1
	}
	return 1
}
*/
