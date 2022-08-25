package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/domain/repository"
)

type SnapshotService interface {
	Create(ctx context.Context, snapshot *model.Snapshot) error
	Update(ctx context.Context, snapshot *model.Snapshot) error
	DeleteByAggregateId(ctx context.Context, tenantId string, aggregateId string) error
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) ([]*model.Snapshot, bool, error)
	FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.Snapshot, bool, error)
}

type snapshotService struct {
	repos repository.SnapshotRepository
}

func NewSnapshotService(repos repository.SnapshotRepository) SnapshotService {
	return &snapshotService{repos: repos}
}

func (s *snapshotService) Create(ctx context.Context, snapshot *model.Snapshot) error {
	if snapshot == nil {
		return nil
	}
	snapshot.TimeStamp = utils.NewMongoNow()
	return s.repos.Create(ctx, snapshot.TenantId, snapshot)
}

func (s *snapshotService) DeleteByAggregateId(ctx context.Context, tenantId, aggregateId string) error {
	return s.repos.DeleteByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) Update(ctx context.Context, snapshot *model.Snapshot) error {
	if snapshot == nil {
		return nil
	}
	return s.repos.Update(ctx, snapshot.TenantId, snapshot)
}

func (s *snapshotService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) ([]*model.Snapshot, bool, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.Snapshot, bool, error) {
	return s.repos.FindByMaxSequenceNumber(ctx, tenantId, aggregateId, aggregateType)
}
