package gorm_impl

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl/db"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/impl/gorm_impl/repository_impl"
)

const (
	ComponentSpecMySql = "eventstorage.mysql"
)

func NewMySqlOptions(log logger.Logger, metadata common.Metadata, adapter eventstorage.GetPubsubAdapter) (*eventstorage.Options, error) {
	mysqlConfig := db.NewMySqlDB(log)
	if err := mysqlConfig.Init(metadata); err != nil {
		return nil, err
	}

	gormDB := mysqlConfig.GetDB()

	config := mysqlConfig.StorageMetadata()
	event := repository_impl.NewEventRepository(gormDB)
	snapshot := repository_impl.NewSnapshotRepository(gormDB)
	aggregate := repository_impl.NewAggregateRepository(gormDB)
	relation := repository_impl.NewRelationRepository(gormDB)
	message := repository_impl.NewMessageRepository(gormDB)
	session := db.NewSession(gormDB)

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
