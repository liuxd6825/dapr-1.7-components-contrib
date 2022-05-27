package service

import (
	"context"
	"github.com/dapr/components-contrib/liuxd/common"
	"github.com/dapr/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/dapr/components-contrib/liuxd/eventstorage/es_mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type SnapshotService interface {
	Create(ctx context.Context, snapshot *model.SnapshotEntity) error
	Update(ctx context.Context, snapshot *model.SnapshotEntity) error
	FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error)
	FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*model.SnapshotEntity, error)
}

func NewSnapshotService(client *mongo.Client, collection *mongo.Collection) SnapshotService {
	return &snapshotService{
		repos: repository.NewSnapshotRepository(client, collection),
	}
}

type snapshotService struct {
	repos *repository.SnapshotRepository
}

func (s *snapshotService) Create(ctx context.Context, snapshot *model.SnapshotEntity) error {
	snapshot.TimeStamp = common.NewMongoNow()
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) Update(ctx context.Context, snapshot *model.SnapshotEntity) error {
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*model.SnapshotEntity, error) {
	return s.repos.FindByMaxSequenceNumber(ctx, tenantId, aggregateId)
}
