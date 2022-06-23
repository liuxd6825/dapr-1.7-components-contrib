package repository

import (
	"context"
	"github.com/liuxd6825/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/model"
	"github.com/liuxd6825/components-contrib/liuxd/eventstorage/es_mongo/other"
	"github.com/orcaman/concurrent-map"
	"go.mongodb.org/mongo-driver/mongo"
)

var collections = cmap.New()

type RelationRepository struct {
	BaseRepository[*model.RelationEntity]
}

func NewRelationRepository(mongodb *other.MongoDB) *RelationRepository {
	res := &RelationRepository{}
	res.mongodb = mongodb
	return res
}

func (r *RelationRepository) InsertOne(ctx context.Context, relation *model.RelationEntity) error {
	coll := r.GetCollection(relation.TableName)
	_, err := coll.InsertOne(ctx, relation)
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationRepository) InsertMany(ctx context.Context, tableName string, relations []*model.RelationEntity) error {
	var docs []interface{}
	for _, rel := range relations {
		docs = append(docs, rel)
	}
	coll := r.GetCollection(tableName)
	_, err := coll.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationRepository) UpdateOne(ctx context.Context, relation *model.RelationEntity) error {
	coll := r.GetCollection(relation.TableName)
	_, err := coll.UpdateByID(ctx, relation.Id, relation)
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationRepository) FindPaging(ctx context.Context, tableName string, query eventstorage.FindPagingQuery, opts ...*other.FindOptions) *eventstorage.FindPagingResult[*model.RelationEntity] {
	coll := r.GetCollection(utils.AsMongoName(tableName))
	return r.BaseRepository.FindPaging(ctx, coll, query, opts...)
}

func (r *RelationRepository) GetCollection(name string) *mongo.Collection {
	value, ok := collections.Get(name)
	if !ok {
		value = r.BaseRepository.mongodb.NewCollection(name)
		collections.Set(name, value)
	}
	coll, _ := value.(*mongo.Collection)
	return coll
}
