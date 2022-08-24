package common

import (
	"errors"
	"fmt"
	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/components-contrib/state"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// MySqlDB is a state store implementation for MongoDB.
type MySqlDB struct {
	state.DefaultBulkStore
	db                *gorm.DB
	operationTimeout  time.Duration
	metadata          MySqlMetadata
	logger            logger.Logger
	collectionOptions *options.CollectionOptions
	snapshotTrigger   uint64
}

type MySqlMetadata struct {
	host             string
	port             string
	username         string
	password         string
	databaseName     string
	params           string
	operationTimeout time.Duration
	maxPoolSize      uint64
}

// NewMySqlDB returns a new MongoDB state store.
func NewMySqlDB(logger logger.Logger) *MySqlDB {
	s := &MySqlDB{
		logger: logger,
	}
	return s
}

// Init establishes connection to the store based on the metadata.
func (m *MySqlDB) Init(metadata Metadata) error {
	meta, err := m.newMySqlMetaData(metadata)
	if err != nil {
		return err
	}

	m.operationTimeout = meta.operationTimeout

	db, err := m.newMySqlDB(meta)
	if err != nil {
		return fmt.Errorf("error in creating mongodb client: %s", err)
	}
	m.db = db
	m.metadata = *meta
	return nil
}

func (m *MySqlDB) GetDB() *gorm.DB {
	return m.db
}

func (m *MySqlDB) GetMetadata() *MySqlMetadata {
	return &m.metadata
}

func (m *MySqlDB) getMySqlURI(metadata *MySqlMetadata) string {
	// root:11111111@tcp(127.0.0.1:3306)/dapr_es?charset=utf8&parseTime=True&loc=Local
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", metadata.username, metadata.password, metadata.host, metadata.port, metadata.databaseName)
}

func (m *MySqlDB) newMySqlDB(metadata *MySqlMetadata) (*gorm.DB, error) {
	uri := m.getMySqlURI(metadata)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       uri,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		m.logger.Error(err.Error())
	}
	return db, err
}

func (m *MySqlDB) newMySqlMetaData(metadata Metadata) (*MySqlMetadata, error) {
	meta := MySqlMetadata{
		databaseName:     defaultDatabaseName,
		operationTimeout: defaultTimeout,
	}

	if val, ok := metadata.Properties[host]; ok && val != "" {
		meta.host = val
	}

	if val, ok := metadata.Properties[port]; ok && val != "" {
		meta.port = val
	}

	if len(meta.host) == 0 {
		return nil, errors.New("must set 'host' fields in metadata")
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

	if val, ok := metadata.Properties[port]; ok && val != "" {
		meta.port = val
	}

	if val, ok := metadata.Properties[maxPoolSize]; ok && val != "" {
		size, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			panic(fmt.Sprintf("%s %s is not uint64", maxPoolSize, val))
		}
		meta.maxPoolSize = size
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
