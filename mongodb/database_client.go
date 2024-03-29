package mongodb

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type databaseClient struct {
	mClient  *mongo.Client
	database *mongo.Database
	logger   *slog.Logger
}

func NewDatabaseClient() DatabaseClient {
	c := &databaseClient{}

	return c

}

func (c *databaseClient) WithDatabase(database string, databaseOptions *options.DatabaseOptions) DatabaseClient {
	c.database = c.mClient.Database(database, databaseOptions)
	return c
}

func (c *databaseClient) WithConnection(mongoUri string, opts *options.ClientOptions) DatabaseClient {
	opts.ApplyURI(mongoUri)

	newClient, err := mongo.NewClient(opts)
	if err != nil {
		c.logger.Error("unable to get mongodb client",
			slog.Any("error", err),
		)
		panic(err)
	}
	if err := newClient.Connect(context.Background()); err != nil {
		log.Println(err)
	}
	c.mClient = newClient

	return c
}

func (c *databaseClient) WithLogger(logger *slog.Logger) DatabaseClient {
	c.logger = logger

	return c
}

func (c *databaseClient) MustValidate() {
	if c.mClient == nil {
		panic("mongo client not set, use WithConnection")
	}

	if c.database == nil {
		panic("database not set, use WithDatabase")

	}

	if c.logger == nil {
		c.logger = slog.Default()

	}

}

func (c *databaseClient) GetDocumentById(ctx context.Context, recordID interface{}, collectionName string) (MongoCursor, error) {
	collection := c.database.Collection(collectionName)
	if collection == nil {
		return nil, errors.New("collection not found")
	}

	result, err := FindDocument(ctx, collection, primitive.M{"_id": recordID})
	if err != nil {
		c.logger.Error("unable to find document",
			slog.Any("error", err),
		)
	}

	if !result.HasData {

		c.logger.Warn("document contains no data",
			slog.Any("results", result),
		)

		return nil, nil
	}

	return result.MongoCursor, nil
}

func (c *databaseClient) SaveDocument(ctx context.Context, collectionName string, document interface{}, modelID interface{}) (bool, error) {
	collection := c.database.Collection(collectionName)
	if collection == nil {
		return false, errors.New("invalid collection")
	}

	res, err := SaveOrUpdateDocument(ctx, collection, document, modelID)
	if err != nil {
		c.logger.Error("insert failed, attempting to replace",
			slog.Any("model_id", modelID),
		)
		return false, err
	}

	return res.IsSuccess, nil

}

func (c *databaseClient) ScanCollection(ctx context.Context, collectionName string) (MongoCursor, error) {
	collection := c.database.Collection(collectionName)
	if collection == nil {
		return nil, errors.New("invalid collection")
	}

	res, err := FindDocument(ctx, collection, primitive.M{})
	if err != nil {
		return nil, err
	}

	if !res.HasData {
		return nil, nil
	}

	return res.MongoCursor, err
}

func (c *databaseClient) CountDocuments(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	collection := c.database.Collection(collectionName)
	if collection == nil {
		return 0, errors.New("invalid collection")
	}

	switch f := filter.(type) {
	case *primitive.M:
		count, err := collection.CountDocuments(ctx, f)
		if err != nil {
			return 0, err
		}
		return count, nil

	default:
		c.logger.Warn("no valid type detected")
		return 0, errors.New("invalid type")

	}

}

func (c *databaseClient) CloseConnection(ctx context.Context) error {
	return c.mClient.Disconnect(ctx)
}

func (c *databaseClient) GetMongoClient() *mongo.Client {
	return c.mClient
}

func (c *databaseClient) GetDatabase() *mongo.Database {
	return c.database
}
