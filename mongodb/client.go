package mongodb

import (
	"context"
	"github.com/byt3-m3/goutils/env_utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type IMongoClient interface {
	GetDocumentById(ctx context.Context, recordId interface{}) (MongoCursor, error)
	SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
	CountDocuments(ctx context.Context, filter primitive.M) (int64, error)
	CloseConnection(ctx context.Context) error
}

type ClientConfig struct {
	MClientConfig  *MongoClientConfig
	CollectionName string
	DBName         string
}

type client struct {
	mClient    *mongo.Client
	collection *mongo.Collection
}

func ProvideIMongoClientConfig(mongoUri, dbName, collectionName string) *ClientConfig {
	return &ClientConfig{
		MClientConfig:  &MongoClientConfig{MongoURI: mongoUri},
		CollectionName: dbName,
		DBName:         collectionName,
	}
}

func ProvideIMongoClientConfigFromEnv(dbName, collectionName string) *ClientConfig {
	return &ClientConfig{
		MClientConfig:  &MongoClientConfig{MongoURI: env_utils.GetEnvStrict("MONGO_URI")},
		CollectionName: dbName,
		DBName:         collectionName,
	}
}

func ProvideIMongoClient(cfg *ClientConfig) IMongoClient {
	mClient := getMongoClientV1(cfg.MClientConfig)
	collection := GetCollectionV1(mClient, cfg.DBName, cfg.CollectionName)

	return &client{
		mClient:    mClient,
		collection: collection,
	}
}

func (c *client) GetDocumentById(ctx context.Context, recordID interface{}) (MongoCursor, error) {
	result, err := FindDocument(ctx, c.collection, primitive.M{"_id": recordID})
	if err != nil {
		log.Println(err)
	}

	if !result.HasData {
		return nil, nil
	}

	return result.MongoCursor, nil
}

func (c *client) SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error) {
	res, err := SaveOrUpdateDocument(ctx, c.collection, document, modelID)
	if err != nil {
		return false, err
	}

	return res.IsSuccess, nil

}

func (c *client) ScanCollection(ctx context.Context) (*mongo.Cursor, error) {
	res, err := FindDocument(ctx, c.collection, primitive.M{})
	if err != nil {
		return nil, err
	}

	if !res.HasData {
		return nil, nil
	}

	return res.MongoCursor, err
}

func (c client) CountDocuments(ctx context.Context, filter primitive.M) (int64, error) {
	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c client) CloseConnection(ctx context.Context) error {
	return c.mClient.Disconnect(ctx)
}
