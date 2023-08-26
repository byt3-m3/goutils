package client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMongoClient interface {
	GetDocumentById(ctx context.Context, recordId interface{}) (*mongo.Cursor, error)
	SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
	CountDocuments(ctx context.Context, filter primitive.M) (int64, error)
	CloseConnection(ctx context.Context) error
	GetMongoCollection() *mongo.Collection
	GetMongoClient() *mongo.Client
}
