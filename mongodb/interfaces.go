package mongodb

import (
	"context"
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
)
