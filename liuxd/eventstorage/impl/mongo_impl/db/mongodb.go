/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	*common.MongoDBMetadata
	aggregateCollectionName string
	eventCollectionName     string
	snapshotCollectionName  string
	relationCollectionName  string
	messageCollectionName   string
	snapshotCount           uint64
}

// MongoDbConfig is a state store implementation for MongoDbConfig.
type MongoDbConfig struct {
	*common.MongoDB
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
func (m *MongoDbConfig) StorageMetadata() *StorageMetadata {
	return m.storageMetadata
}

// NewMongoDB returns a new MongoDbConfig state store.
func NewMongoDB(logger logger.Logger) *MongoDbConfig {
	mdb := common.NewMongoDB(logger)
	s := &MongoDbConfig{
		MongoDB: mdb,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (m *MongoDbConfig) Init(metadata common.Metadata) error {
	if err := m.MongoDB.Init(metadata); err != nil {
		return err
	}
	storageMetadata, err := m.getStorageMetadata(metadata)
	if err != nil {
		return err
	}
	m.storageMetadata = storageMetadata
	return nil
}

func (m *MongoDbConfig) getStorageMetadata(metadata common.Metadata) (*StorageMetadata, error) {
	meta := StorageMetadata{
		MongoDBMetadata:         m.MongoDB.GetMetadata(),
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
