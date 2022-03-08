package service

import (
	"context"
	"github.com/dapr/components-contrib/eventsourcing/utils"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/model"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type SnapshotService interface {
	Save(ctx context.Context, snapshot *model.SnapshotEntity) error
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

func (s *snapshotService) Save(ctx context.Context, snapshot *model.SnapshotEntity) error {
	snapshot.TimeStamp = utils.NewMongoNow()
	return s.repos.Insert(ctx, snapshot)
}

func (s *snapshotService) FindByAggregateId(ctx context.Context, tenantId string, aggregateId string) (*[]model.SnapshotEntity, error) {
	return s.repos.FindByAggregateId(ctx, tenantId, aggregateId)
}

func (s *snapshotService) FindByMaxSequenceNumber(ctx context.Context, tenantId string, aggregateId string) (*model.SnapshotEntity, error) {
	return s.repos.FindByMaxSequenceNumber(ctx, tenantId, aggregateId)
}
