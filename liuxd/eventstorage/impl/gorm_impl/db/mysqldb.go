package db

import (
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/liuxd/common"
	"strconv"
)

const (
	eventCollectionName     = "eventCollectionName"
	snapshotCollectionName  = "snapshotCollectionName"
	aggregateCollectionName = "aggregateCollectionName"
	relationCollectionName  = "relationCollectionName"
	messageCollectionName   = "messageCollectionName"
	snapshotCountName       = "snapshotCountName"
	id                      = "_id"
	value                   = "value"
	etag                    = "_etag"

	defaultEventCollectionName     = "ddd_event"
	defaultSnapshotCollectionName  = "ddd_snapshot"
	defaultAggregateCollectionName = "ddd_aggregate"
	defaultRelationCollectionName  = "ddd_relation"
	defaultMessageCollectionName   = "ddd_message"
	defaultSnapshotCount           = 100
)

type StorageMetadata struct {
	*common.MySqlMetadata
	aggregateCollectionName string
	eventCollectionName     string
	snapshotCollectionName  string
	relationCollectionName  string
	messageCollectionName   string
	snapshotCount           uint64
}

// MySqlConfig is a state store implementation for MongoDbConfig.
type MySqlConfig struct {
	*common.MySqlDB
	storageMetadata *StorageMetadata
}

func (s *StorageMetadata) AggregateCollectionName() string {
	return s.aggregateCollectionName
}

func (s *StorageMetadata) EventCollectionName() string {
	return s.eventCollectionName
}

func (s *StorageMetadata) SnapshotCollectionName() string {
	return s.snapshotCollectionName
}

func (s *StorageMetadata) RelationCollectionName() string {
	return s.relationCollectionName
}

func (s *StorageMetadata) MessageCollectionName() string {
	return s.messageCollectionName
}
func (s *StorageMetadata) SnapshotCount() uint64 {
	return s.snapshotCount
}

func (s *MySqlConfig) StorageMetadata() *StorageMetadata {
	return s.storageMetadata
}

// NewMySqlDB returns a new MongoDbConfig state store.
func NewMySqlDB(logger logger.Logger) *MySqlConfig {
	mdb := common.NewMySqlDB(logger)
	s := &MySqlConfig{
		MySqlDB: mdb,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (s *MySqlConfig) Init(metadata common.Metadata) error {
	if err := s.MySqlDB.Init(metadata); err != nil {
		return err
	}
	storageMetadata, err := s.getStorageMetadata(metadata)
	if err != nil {
		return err
	}
	s.storageMetadata = storageMetadata
	return nil
}

func (s *MySqlConfig) getStorageMetadata(metadata common.Metadata) (*StorageMetadata, error) {
	meta := StorageMetadata{
		MySqlMetadata:           s.MySqlDB.GetMetadata(),
		eventCollectionName:     defaultEventCollectionName,
		snapshotCollectionName:  defaultSnapshotCollectionName,
		aggregateCollectionName: defaultAggregateCollectionName,
		relationCollectionName:  defaultRelationCollectionName,
		messageCollectionName:   defaultMessageCollectionName,
		snapshotCount:           defaultSnapshotCount,
	}
	if val, ok := metadata.Properties[eventCollectionName]; ok && val != "" {
		meta.eventCollectionName = val
	}
	if val, ok := metadata.Properties[snapshotCollectionName]; ok && val != "" {
		meta.snapshotCollectionName = val
	}
	if val, ok := metadata.Properties[aggregateCollectionName]; ok && val != "" {
		meta.aggregateCollectionName = val
	}
	if val, ok := metadata.Properties[relationCollectionName]; ok && val != "" {
		meta.relationCollectionName = val
	}
	if val, ok := metadata.Properties[messageCollectionName]; ok && val != "" {
		meta.messageCollectionName = val
	}
	if val, ok := metadata.Properties[snapshotCountName]; ok && val != "" {
		count, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		meta.snapshotCount = count
	}
	return &meta, nil
}
