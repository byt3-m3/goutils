package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type collectionClient struct {
	mClient    *mongo.Client
	collection *mongo.Collection
	logger     *slog.Logger
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
		c.logger.Error("unable to create new mongodb client",
			slog.Any("error", err),
			slog.Any("opts", opts),
		)
	}
	if err := newClient.Connect(context.Background()); err != nil {
		c.logger.Error("unable to connect to server",
			slog.Any("error", err),
			slog.Any("opts", opts),
		)
	}
	c.mClient = newClient

	return c
}

func (c *collectionClient) WithLogger(logger *slog.Logger) CollectionClient {
	c.logger = logger

	return c
}

func (c *collectionClient) MustValidate() {
	if c.mClient == nil {
		panic("mongo client not set, use WithConnection")
	}

	if c.collection == nil {
		panic("collection not set, use WithCollection")

	}

	if c.logger == nil {
		c.logger = slog.Default()

	}

}

func (c *collectionClient) GetDocumentById(ctx context.Context, recordID interface{}) (MongoCursor, error) {
	result, err := FindDocument(ctx, c.collection, primitive.M{"_id": recordID})
	if err != nil {
		c.logger.Error("unable to get document by id",
			slog.Any("error", err),
			slog.Any("record_id", recordID),
		)
		return nil, err
	}

	if !result.HasData {
		c.logger.Warn("document contains no data",
			slog.Any("results", result),
		)

		return nil, nil
	}

	return result.MongoCursor, nil
}

func (c *collectionClient) SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error) {
	res, err := SaveOrUpdateDocument(ctx, c.collection, document, modelID)
	if err != nil {
		c.logger.Error("insert failed, attempting to replace: %s",
			slog.Any("model_id", modelID),
		)

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
		c.logger.Warn("no valid type detected")
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

func (c *collectionClient) DeleteDocument(ctx context.Context, filter interface{}, logger slog.Logger) error {
	result, err := c.mClient.Database(c.collection.Database().Name()).Collection(c.collection.Name()).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	logger.Debug("deleted records",
		slog.Int("count", int(result.DeletedCount)))
	return nil
}
