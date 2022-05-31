package common

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/state"
)

const (
	host             = "host"
	username         = "username"
	password         = "password"
	databaseName     = "databaseName"
	server           = "server"
	writeConcern     = "writeConcern"
	readConcern      = "readConcern"
	operationTimeout = "operationTimeout"
	params           = "params"
	replicaSet       = "replica-set"
	maxPoolSize      = "max-pool-size"
	id               = "_id"
	value            = "value"
	etag             = "_etag"

	defaultTimeout      = 5 * time.Second
	defaultDatabaseName = ""

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
	client            *mongo.Client
	operationTimeout  time.Duration
	database          *mongo.Database
	metadata          MongoDBMetadata
	logger            logger.Logger
	collectionOptions *options.CollectionOptions
}

type MongoDBMetadata struct {
	host             string
	username         string
	password         string
	databaseName     string
	server           string
	writeconcern     string
	readconcern      string
	params           string
	operationTimeout time.Duration
	replicaSet       string
	maxPoolSize      uint64
}

// NewMongoDB returns a new MongoDB state store.
func NewMongoDB(logger logger.Logger) *MongoDB {
	s := &MongoDB{
		logger: logger,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (m *MongoDB) Init(metadata Metadata) error {
	meta, err := m.getMongoDBMetaData(metadata)
	if err != nil {
		return err
	}

	m.operationTimeout = meta.operationTimeout

	client, err := m.getMongoDBClient(meta)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("error in connecting to mongodb, host: %s error: %s", meta.host, err)
	}

	m.client = client

	// get the write concern
	wc, err := m.getWriteConcernObject(meta.writeconcern)
	if err != nil {
		return fmt.Errorf("error in getting write concern object: %s", err)
	}

	// get the read concern
	rc, err := m.getReadConcernObject(meta.readconcern)
	if err != nil {
		return fmt.Errorf("error in getting read concern object: %s", err)
	}

	m.metadata = *meta
	m.collectionOptions = options.Collection().SetWriteConcern(wc).SetReadConcern(rc)
	m.database = m.client.Database(meta.databaseName)

	return nil
}

func (m *MongoDB) GetClient() *mongo.Client {
	return m.client
}

func (m *MongoDB) NewCollection(name string) *mongo.Collection {
	return m.database.Collection(name, m.collectionOptions)
}

func (m *MongoDB) Ping() error {
	if err := m.client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("mongoDB store: error connecting to mongoDB at %s: %s", m.metadata.host, err)
	}

	return nil
}

func (m *MongoDB) GetMetadata() *MongoDBMetadata {
	return &m.metadata
}

func (m *MongoDB) getMongoURI(metadata *MongoDBMetadata) string {
	if len(metadata.server) != 0 {
		return fmt.Sprintf(connectionURIFormatWithSrv, metadata.server, metadata.params)
	}

	if metadata.username != "" && metadata.password != "" {
		return fmt.Sprintf(connectionURIFormatWithAuthentication, metadata.username, metadata.password, metadata.host, metadata.databaseName, metadata.params)
	}

	return fmt.Sprintf(connectionURIFormat, metadata.host, metadata.databaseName, metadata.params)
}

func (m *MongoDB) getMongoDBClient(metadata *MongoDBMetadata) (*mongo.Client, error) {
	uri := m.getMongoURI(metadata)

	// Set client options
	opts := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), metadata.operationTimeout)
	defer cancel()

	daprUserAgent := "dapr-event-sourcing-" + logger.DaprVersion
	if opts.AppName != nil {
		opts.SetAppName(daprUserAgent + ":" + *opts.AppName)
	} else {
		opts.SetAppName(daprUserAgent)
	}

	if metadata.replicaSet != "" {
		opts.SetReplicaSet(metadata.replicaSet)
	}

	if metadata.maxPoolSize != 0 {
		opts.SetMaxPoolSize(metadata.maxPoolSize)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MongoDB) getMongoDBMetaData(metadata Metadata) (*MongoDBMetadata, error) {
	meta := MongoDBMetadata{
		databaseName:     defaultDatabaseName,
		operationTimeout: defaultTimeout,
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

	if val, ok := metadata.Properties[writeConcern]; ok && val != "" {
		meta.writeconcern = val
	}

	if val, ok := metadata.Properties[readConcern]; ok && val != "" {
		meta.readconcern = val
	}

	if val, ok := metadata.Properties[params]; ok && val != "" {
		meta.params = val
	}

	if val, ok := metadata.Properties[replicaSet]; ok && val != "" {
		meta.replicaSet = val
	}

	if val, ok := metadata.Properties[maxPoolSize]; ok && val != "" {
		size, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			meta.maxPoolSize = size
		} else {
			panic(fmt.Sprintf("%s %s is not nint64", maxPoolSize, val))
		}
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

func (m *MongoDB) getWriteConcernObject(cn string) (*writeconcern.WriteConcern, error) {
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

func (m *MongoDB) getReadConcernObject(cn string) (*readconcern.ReadConcern, error) {
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
