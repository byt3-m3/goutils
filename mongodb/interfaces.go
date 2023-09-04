package mongodb

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	MongoCursor interface {
		Decode(val interface{}) error
		Err() error
		Next(ctx context.Context) bool
		Close(ctx context.Context) error
		ID() int64
	}

	MongoCollectionReplacerInserter interface {
		MongoCollectionReplacer
		MongoCollectionInserter
	}
	MongoCollectionFinderReplacerInserter interface {
		MongoCollectionReplacer
		MongoCollectionInserter
		MongoCollectionFinder
	}
	MongoCollection interface {
		MongoCollectionReplacer
		MongoCollectionInserter
		MongoCollectionFinder
		MongoCollectionCounter
	}

	MongoCollectionReplacer interface {
		ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)
	}

	MongoCollectionInserter interface {
		InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	}

	MongoCollectionFinder interface {
		Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	}

	MongoCollectionCounter interface {
		CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	}

	MongoDocumentSaver interface {
		SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
	}

	MongoDocumentQuerier interface {
		GetDocumentById(ctx context.Context, recordId interface{}) (MongoCursor, error)
	}

	MongoCollectionGetter interface {
		GetMongoCollection() *mongo.Collection
	}

	MongoConnectionCloser interface {
		CloseConnection(ctx context.Context) error
	}

	MongoCollectionScanner interface {
		ScanCollection(ctx context.Context) (MongoCursor, error)
	}

	MongoClient interface {
		MongoDocumentSaver
		MongoCollectionCounter
		MongoConnectionCloser
		MongoCollectionGetter
		GetMongoClient() *mongo.Client
		MongoCollectionScanner
	}

	DatabaseClientOptionSetter interface {
		WithDatabase(database string, databaseOptions *options.DatabaseOptions) DatabaseClient
		WithConnection(mongoUri string, opts *options.ClientOptions) DatabaseClient

		WithLogger(logger *log.Logger) DatabaseClient
	}

	DatabaseClient interface {
		DatabaseClientOptionSetter
		GetDocumentById(ctx context.Context, recordID interface{}, collectionName string) (MongoCursor, error)
		SaveDocument(ctx context.Context, collectionName string, document interface{}, modelID interface{}) (bool, error)
		ScanCollection(ctx context.Context, collectionName string) (MongoCursor, error)
		CountDocuments(ctx context.Context, collectionName string, filter interface{}) (int64, error)
		CloseConnection(ctx context.Context) error
		GetMongoClient() *mongo.Client
		GetDatabase() *mongo.Database
		MustValidate()
	}

	CollectionClientOptionSetter interface {
		WithCollection(database, collection string, databaseOptions *options.DatabaseOptions) CollectionClient
		WithConnection(mongoUri string, opts *options.ClientOptions) CollectionClient
		WithLogger(logger *log.Logger) CollectionClient
	}

	CollectionClient interface {
		CollectionClientOptionSetter

		GetDocumentById(ctx context.Context, recordID interface{}) (MongoCursor, error)
		SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
		ScanCollection(ctx context.Context) (MongoCursor, error)
		CountDocuments(ctx context.Context, filter interface{}) (int64, error)
		CloseConnection(ctx context.Context) error
		GetMongoClient() *mongo.Client
		GetCollection() *mongo.Collection
		MustValidate()
	}
)
