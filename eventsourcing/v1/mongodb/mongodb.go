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

package mongodb

// mongodb package is an implementation of StateStore interface to perform operations on store

import (
	"context"
	"errors"
	"fmt"
	es "github.com/dapr/components-contrib/eventsourcing/v1"
	"github.com/dapr/components-contrib/eventsourcing/v1/mongodb/domain/service"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/kit/logger"
)

const (
	host                   = "host"
	username               = "username"
	password               = "password"
	databaseName           = "databaseName"
	eventCollectionName    = "eventCollectionName"
	snapshotCollectionName = "snapshotCollectionName"
	server                 = "server"
	writeConcern           = "writeConcern"
	readConcern            = "readConcern"
	operationTimeout       = "operationTimeout"
	params                 = "params"
	id                     = "_id"
	value                  = "value"
	etag                   = "_etag"

	defaultTimeout                = 5 * time.Second
	defaultDatabaseName           = "dapr_esdb"
	defaultEventCollectionName    = "dapr_event"
	defaultSnapshotCollectionName = "dapr_snapshot"

	// mongodb://<username>:<password@<host>/<database><params>
	connectionURIFormatWithAuthentication = "mongodb://%s:%s@%s/%s%s"

	// mongodb://<host>/<database><params>
	connectionURIFormat = "mongodb://%s/%s%s"

	// mongodb+srv://<server>/<params>
	connectionURIFormatWithSrv = "mongodb+srv://%s/%s"
)

// MongoDB is a state store implementation for MongoDB.
type MongoDB struct {
	state.DefaultBulkStore
	client           *mongo.Client
	operationTimeout time.Duration
	metadata         mongoDBMetadata

	logger logger.Logger

	eventService     service.EventService
	snapshotService  service.SnapshotService
	aggregateService service.AggregateService
	eventLogService  service.EventLogService
}

type mongoDBMetadata struct {
	host                   string
	username               string
	password               string
	databaseName           string
	eventCollectionName    string
	snapshotCollectionName string
	server                 string
	writeconcern           string
	readconcern            string
	params                 string
	operationTimeout       time.Duration
}

// NewMongoDB returns a new MongoDB state store.
func NewMongoDB(logger logger.Logger) *MongoDB {
	s := &MongoDB{
		logger: logger,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (m *MongoDB) Init(metadata es.Metadata) error {
	meta, err := getMongoDBMetaData(metadata)
	if err != nil {
		return err
	}

	m.operationTimeout = meta.operationTimeout

	client, err := getMongoDBClient(meta)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("error in connecting to mongodb, host: %s error: %s", meta.host, err)
	}

	m.client = client

	// get the write concern
	wc, err := getWriteConcernObject(meta.writeconcern)
	if err != nil {
		return fmt.Errorf("error in getting write concern object: %s", err)
	}

	// get the read concern
	rc, err := getReadConcernObject(meta.readconcern)
	if err != nil {
		return fmt.Errorf("error in getting read concern object: %s", err)
	}

	m.metadata = *meta
	opts := options.Collection().SetWriteConcern(wc).SetReadConcern(rc)

	database := m.client.Database(meta.databaseName)

	eventCollection := database.Collection(meta.eventCollectionName, opts)
	snapshotCollection := database.Collection(meta.snapshotCollectionName, opts)

	m.eventService = service.NewEventService(m.client, eventCollection)
	m.snapshotService = service.NewSnapshotService(m.client, snapshotCollection)
	m.aggregateService = service.NewAggregateService(m.client, eventCollection)
	return nil
}

func (m *MongoDB) Ping() error {
	if err := m.client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("mongoDB store: error connecting to mongoDB at %s: %s", m.metadata.host, err)
	}

	return nil
}

func getMongoURI(metadata *mongoDBMetadata) string {
	if len(metadata.server) != 0 {
		return fmt.Sprintf(connectionURIFormatWithSrv, metadata.server, metadata.params)
	}

	if metadata.username != "" && metadata.password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, metadata.username, metadata.password, metadata.host, metadata.databaseName, metadata.params)
	}

	return fmt.Sprintf(connectionURIFormat, metadata.host, metadata.databaseName, metadata.params)
}

func getMongoDBClient(metadata *mongoDBMetadata) (*mongo.Client, error) {
	uri := getMongoURI(metadata)

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), metadata.operationTimeout)
	defer cancel()

	daprUserAgent := "dapr-event-sourcing-" + logger.DaprVersion
	if clientOptions.AppName != nil {
		clientOptions.SetAppName(daprUserAgent + ":" + *clientOptions.AppName)
	} else {
		clientOptions.SetAppName(daprUserAgent)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getMongoDBMetaData(metadata es.Metadata) (*mongoDBMetadata, error) {
	meta := mongoDBMetadata{
		databaseName:           defaultDatabaseName,
		eventCollectionName:    defaultEventCollectionName,
		snapshotCollectionName: defaultSnapshotCollectionName,
		operationTimeout:       defaultTimeout,
	}

	if val, ok := metadata.Properties[host]; ok && val != "" {
		meta.host = val
	}

	if val, ok := metadata.Properties[server]; ok && val != "" {
		meta.server = val
	}

	if len(meta.host) == 0 && len(meta.server) == 0 {
		return nil, errors.New("must set 'host' or 'server' fields in metadata")
	}

	if len(meta.host) != 0 && len(meta.server) != 0 {
		return nil, errors.New("'host' or 'server' fields are mutually exclusive")
	}

	if val, ok := metadata.Properties[username]; ok && val != "" {
		meta.username = val
	}

	if val, ok := metadata.Properties[password]; ok && val != "" {
		meta.password = val
	}

	if val, ok := metadata.Properties[databaseName]; ok && val != "" {
		meta.databaseName = val
	}

	if val, ok := metadata.Properties[eventCollectionName]; ok && val != "" {
		meta.eventCollectionName = val
	}

	if val, ok := metadata.Properties[snapshotCollectionName]; ok && val != "" {
		meta.snapshotCollectionName = val
	}

	if val, ok := metadata.Properties[writeConcern]; ok && val != "" {
		meta.writeconcern = val
	}

	if val, ok := metadata.Properties[readConcern]; ok && val != "" {
		meta.readconcern = val
	}

	if val, ok := metadata.Properties[params]; ok && val != "" {
		meta.params = val
	}

	var err error
	if val, ok := metadata.Properties[operationTimeout]; ok && val != "" {
		meta.operationTimeout, err = time.ParseDuration(val)
		if err != nil {
			return nil, errors.New("incorrect operationTimeout field from metadata")
		}
	}

	return &meta, nil
}

func getWriteConcernObject(cn string) (*writeconcern.WriteConcern, error) {
	var wc *writeconcern.WriteConcern
	if cn != "" {
		if cn == "majority" {
			wc = writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
		} else {
			w, err := strconv.Atoi(cn)
			wc = writeconcern.New(writeconcern.W(w), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))

			return wc, err
		}
	} else {
		wc = writeconcern.New(writeconcern.W(1), writeconcern.J(true), writeconcern.WTimeout(defaultTimeout))
	}

	return wc, nil
}

func getReadConcernObject(cn string) (*readconcern.ReadConcern, error) {
	switch cn {
	case "local":
		return readconcern.Local(), nil
	case "majority":
		return readconcern.Majority(), nil
	case "available":
		return readconcern.Available(), nil
	case "linearizable":
		return readconcern.Linearizable(), nil
	case "snapshot":
		return readconcern.Snapshot(), nil
	case "":
		return readconcern.Local(), nil
	}

	return nil, fmt.Errorf("readConcern %s not found", cn)
}
