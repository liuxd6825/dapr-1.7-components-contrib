package service

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type SnapshotService interface {
	Create(ctx context.Context, snapshot *model.SnapshotEntity) error
	Update(ctx context.Context, snapshot *model.SnapshotEntity) error
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error)
	FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.SnapshotEntity, error)
}

func NewSnapshotService(mongodb *db.MongoDB, collection *mongo.Collection) SnapshotService {
	return &snapshotService{
		repos: repository.NewSnapshotRepository(mongodb, collection),
	}
}

type snapshotService struct {
	repos *repository.SnapshotRepository
}

func (s *snapshotService) Create(ctx context.Context, snapshot *model.SnapshotEntity) error {
	snapshot.TimeStamp = utils.NewMongoNow()
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) Update(ctx context.Context, snapshot *model.SnapshotEntity) error {
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string, aggregateType string) (*model.SnapshotEntity, error) {
	return s.repos.FindByMaxSequenceNumber(ctx, tenantId, aggregateId, aggregateType)
}
