package es_mongo

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/repository_impl"
)

const (
	ComponentSpecType = "eventstorage.mongodb"
)

func NewOptions(log logger.Logger, metadata common.Metadata, adapter eventstorage.GetPubsubAdapter) (*eventstorage.Options, error) {
	mongoConfig := db.NewMongoDB(log)
	if err := mongoConfig.Init(metadata); err != nil {
		return nil, err
	}
	config := mongoConfig.StorageMetadata()
	event := repository_impl.NewEventRepository(mongoConfig, config.EventCollectionName())
	snapshot := repository_impl.NewSnapshotRepository(mongoConfig, config.SnapshotCollectionName())
	aggregate := repository_impl.NewAggregateRepository(mongoConfig, config.AggregateCollectionName())
	relation := repository_impl.NewRelationRepository(mongoConfig, config.RelationCollectionName())
	message := repository_impl.NewMessageRepository(mongoConfig, config.MessageCollectionName())
	session := db.NewSession(mongoConfig.GetClient())

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
