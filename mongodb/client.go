package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type clientOpt func(c *client)
type WithMongoClientInput struct {
}

var (
	WithCollection = func(dbName, collection string) clientOpt {
		return func(c *client) {
			c.collection = c.mClient.Database(dbName).Collection(collection)
		}
	}

	WithMongoClient = func(mongoUri string, opts *options.ClientOptions) clientOpt {
		return func(c *client) {
			opts.ApplyURI(mongoUri)

			newClient, err := mongo.NewClient(opts)
			if err != nil {
				log.Fatalln(err)
			}
			if err := newClient.Connect(context.Background()); err != nil {
				log.Println(err)
			}
			c.mClient = newClient
		}
	}
)

type MongoClient interface {
	GetDocumentById(ctx context.Context, recordId interface{}) (MongoCursor, error)
	SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
	CountDocuments(ctx context.Context, filter primitive.M) (int64, error)
	CloseConnection(ctx context.Context) error
	GetMongoCollection() *mongo.Collection
	GetMongoClient() *mongo.Client
	ScanCollection(ctx context.Context) (MongoCursor, error)
}

type client struct {
	mClient    *mongo.Client
	collection *mongo.Collection
}

func NewClient(opts ...clientOpt) MongoClient {
	c := &client{}

	for _, opt := range opts {
		opt(c)
	}
	if !validateClient(c) {
		log.Fatalln("failed client validation")
	}

	return c

}

func validateClient(c *client) bool {
	if c.collection == nil {
		log.Println("collection not set, use WithCollection")
		return false
	}

	if c.mClient == nil {
		log.Println("collection not set, use WithCollection")
		return false
	}

	return true
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

func (c *client) ScanCollection(ctx context.Context) (MongoCursor, error) {
	res, err := FindDocument(ctx, c.collection, primitive.M{})
	if err != nil {
		return nil, err
	}

	if !res.HasData {
		return nil, nil
	}

	return res.MongoCursor, err
}

func (c *client) CountDocuments(ctx context.Context, filter primitive.M) (int64, error) {
	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *client) CloseConnection(ctx context.Context) error {
	return c.mClient.Disconnect(ctx)
}

func (c *client) GetMongoClient() *mongo.Client {
	return c.mClient
}

func (c *client) GetMongoCollection() *mongo.Collection {
	return c.collection
}
