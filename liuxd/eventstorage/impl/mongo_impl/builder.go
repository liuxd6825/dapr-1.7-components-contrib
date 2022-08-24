package mongo_impl

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	db2 "github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl/db"
	repository_impl2 "github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/mongo_impl/repository_impl"
)

const (
	ComponentSpecMongo = "eventstorage.mongodb"
)

func NewMongoOptions(log logger.Logger, metadata common.Metadata, adapter eventstorage.GetPubsubAdapter) (*eventstorage.Options, error) {
	mongoConfig := db2.NewMongoDB(log)
	if err := mongoConfig.Init(metadata); err != nil {
		return nil, err
	}
	config := mongoConfig.StorageMetadata()
	event := repository_impl2.NewEventRepository(mongoConfig, config.EventCollectionName())
	snapshot := repository_impl2.NewSnapshotRepository(mongoConfig, config.SnapshotCollectionName())
	aggregate := repository_impl2.NewAggregateRepository(mongoConfig, config.AggregateCollectionName())
	relation := repository_impl2.NewRelationRepository(mongoConfig, config.RelationCollectionName())
	message := repository_impl2.NewMessageRepository(mongoConfig, config.MessageCollectionName())
	session := db2.NewSession(mongoConfig.GetClient())

	ops := &eventstorage.Options{
		Metadata:       metadata,
		PubsubAdapter:  adapter,
		EventRepos:     event,
		SnapshotRepos:  snapshot,
		AggregateRepos: aggregate,
		RelationRepos:  relation,
		MessageRepos:   message,
		SnapshotCount:  config.SnapshotCount(),
		Session:        session,
	}
	return ops, nil
}
