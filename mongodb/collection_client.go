package mongodb

import (
	"context"
	"errors"
	"github.com/byt3-m3/goutils/logging"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type collectionClient struct {
	mClient    *mongo.Client
	collection *mongo.Collection
	logger     *log.Logger
}

func NewCollectionClient() CollectionClient {
	c := &collectionClient{}

	return c

}

func (c *collectionClient) WithCollection(database, collection string, databaseOptions *options.DatabaseOptions) CollectionClient {
	c.collection = c.mClient.Database(database, databaseOptions).Collection(collection)
	return c
}

func (c *collectionClient) WithConnection(mongoUri string, opts *options.ClientOptions) CollectionClient {
	opts.ApplyURI(mongoUri)

	newClient, err := mongo.NewClient(opts)
	if err != nil {
		c.logger.Fatalln(err)
	}
	if err := newClient.Connect(context.Background()); err != nil {
		log.Println(err)
	}
	c.mClient = newClient

	return c
}

func (c *collectionClient) WithLogger(logger *log.Logger) CollectionClient {
	c.logger = logger

	return c
}

func (c *collectionClient) MustValidate() {
	switch {
	case c.mClient == nil:
		panic("mongo client not set, use WithConnection")

	case c.collection == nil:
		panic("database not set, use WithCollection")

	case c.logger == nil:
		c.logger = logging.NewLogger()

	}
}

func (c *collectionClient) GetDocumentById(ctx context.Context, recordID interface{}) (MongoCursor, error) {
	result, err := FindDocument(ctx, c.collection, primitive.M{"_id": recordID})
	if err != nil {
		c.logger.Error(err)
	}

	if !result.HasData {
		c.logger.WithFields(map[string]interface{}{
			"results": result,
		}).Warning("document contains no data")
		return nil, nil
	}

	return result.MongoCursor, nil
}

func (c *collectionClient) SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error) {
	res, err := SaveOrUpdateDocument(ctx, c.collection, document, modelID)
	if err != nil {
		return false, err
	}

	return res.IsSuccess, nil

}

func (c *collectionClient) ScanCollection(ctx context.Context) (MongoCursor, error) {
	res, err := FindDocument(ctx, c.collection, primitive.M{})
	if err != nil {
		return nil, err
	}

	if !res.HasData {
		return nil, nil
	}

	return res.MongoCursor, err
}

func (c *collectionClient) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {

	switch f := filter.(type) {
	case *primitive.M:
		count, err := c.collection.CountDocuments(ctx, f)
		if err != nil {
			return 0, err
		}
		return count, nil

	default:
		c.logger.Warning("no valid type detected")
		return 0, errors.New("invalid type")

	}

}

func (c *collectionClient) CloseConnection(ctx context.Context) error {
	return c.mClient.Disconnect(ctx)
}

func (c *collectionClient) GetMongoClient() *mongo.Client {
	return c.mClient
}

func (c *collectionClient) GetCollection() *mongo.Collection {
	return c.collection
}
