package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}
