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

package mongo

// mongodb package is an implementation of StateStore interface to perform operations on store

import (
	"github.com/dapr/components-contrib/liuxd/common"
	"github.com/dapr/kit/logger"
)

const (
	appLogCollectionName   = "appLogCollectionName"
	eventLogCollectionName = "eventLogCollectionName"
	id                     = "_id"
	value                  = "value"
	etag                   = "_etag"

	defaultEventLogCollectionName = "dapr_event_logs"
	defaultAppLogCollectionName   = "dapr_app_logs"
)

// MongoDB is a state store implementation for MongoDB.
type MongoDB struct {
	*common.MongoDB
	loggerMetadata *loggerMetadata
}

type loggerMetadata struct {
	*common.MongoDBMetadata
	appLogCollectionName   string
	eventLogCollectionName string
}

// NewMongoDB returns a new MongoDB state store.
func NewMongoDB(logger logger.Logger) *MongoDB {
	mdb := common.NewMongoDB(logger)
	s := &MongoDB{
		MongoDB: mdb,
	}
	return s
}

// Init establishes connection to the store based on the loggerMetadata.
func (m *MongoDB) Init(metadata common.Metadata) error {
	if err := m.MongoDB.Init(metadata); err != nil {
		return err
	}
	loggerMetadata, err := m.getLoggerMetadata(metadata)
	if err != nil {
		return err
	}
	m.loggerMetadata = loggerMetadata
	return nil
}

func (m *MongoDB) getLoggerMetadata(metadata common.Metadata) (*loggerMetadata, error) {
	meta := loggerMetadata{
		MongoDBMetadata:        m.MongoDB.GetMetadata(),
		appLogCollectionName:   defaultAppLogCollectionName,
		eventLogCollectionName: defaultEventLogCollectionName,
	}
	if val, ok := metadata.Properties[appLogCollectionName]; ok && val != "" {
		meta.appLogCollectionName = val
	}
	if val, ok := metadata.Properties[eventLogCollectionName]; ok && val != "" {
		meta.eventLogCollectionName = val
	}
	return &meta, nil
}
