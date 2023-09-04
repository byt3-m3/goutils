package mongodb

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StubDatabaseClient struct {
	GetDocumentByIdReturn func(ctx context.Context, recordID interface{}, collectionName string) (MongoCursor, error)
	SaveDocumentReturn    func(ctx context.Context, collectionName string, document interface{}, modelID interface{}) (bool, error)
	ScanCollectionReturn  func(ctx context.Context, collectionName string) (MongoCursor, error)
	CountDocumentsReturn  func(ctx context.Context, collectionName string, filter interface{}) (int64, error)
	CloseConnectionReturn func(ctx context.Context) error
	GetMongoClientReturn  func() *mongo.Client
	GetDatabaseReturn     func() *mongo.Database
	MustValidateReturn    func()
	WithDatabaseReturn    func(database string, databaseOptions *options.DatabaseOptions)
	WithConnectionReturn  func(mongoUri string, opts *options.ClientOptions)
	WithLoggerReturn      func(logger *log.Logger)
}

func (s StubDatabaseClient) WithDatabase(database string, databaseOptions *options.DatabaseOptions) DatabaseClient {
	s.WithDatabaseReturn(database, databaseOptions)
	return s
}

func (s StubDatabaseClient) WithConnection(mongoUri string, opts *options.ClientOptions) DatabaseClient {
	s.WithConnectionReturn(mongoUri, opts)
	return s
}

func (s StubDatabaseClient) WithLogger(logger *log.Logger) DatabaseClient {
	s.WithLoggerReturn(logger)
	return s
}

func (s StubDatabaseClient) GetDocumentById(ctx context.Context, recordID interface{}, collectionName string) (MongoCursor, error) {
	return s.GetDocumentByIdReturn(ctx, recordID, collectionName)
}

func (s StubDatabaseClient) SaveDocument(ctx context.Context, collectionName string, document interface{}, modelID interface{}) (bool, error) {
	return s.SaveDocumentReturn(ctx, collectionName, document, modelID)
}

func (s StubDatabaseClient) ScanCollection(ctx context.Context, collectionName string) (MongoCursor, error) {
	return s.ScanCollectionReturn(ctx, collectionName)
}

func (s StubDatabaseClient) CountDocuments(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	return s.CountDocumentsReturn(ctx, collectionName, filter)
}

func (s StubDatabaseClient) CloseConnection(ctx context.Context) error {
	return s.CloseConnectionReturn(ctx)
}

func (s StubDatabaseClient) GetMongoClient() *mongo.Client {
	return s.GetMongoClientReturn()
}

func (s StubDatabaseClient) GetDatabase() *mongo.Database {
	return s.GetDatabaseReturn()
}

func (s StubDatabaseClient) MustValidate() {
	s.MustValidateReturn()
}

type StubCollectionClient struct {
	GetDocumentByIdReturn func(ctx context.Context, recordID interface{}) (MongoCursor, error)
	SaveDocumentReturn    func(ctx context.Context, document interface{}, modelID interface{}) (bool, error)
	ScanCollectionReturn  func(ctx context.Context) (MongoCursor, error)
	CountDocumentsReturn  func(ctx context.Context, filter interface{}) (int64, error)
	CloseConnectionReturn func(ctx context.Context) error
	GetMongoClientReturn  func() *mongo.Client
	GetCollectionReturn   func() *mongo.Collection
	MustValidateReturn    func()
	WithCollectionReturn  func(database, collection string, databaseOptions *options.DatabaseOptions)
	WithConnectionReturn  func(mongoUri string, opts *options.ClientOptions)
	WithLoggerReturn      func(logger *log.Logger)
}

func (s StubCollectionClient) WithCollection(database, collection string, databaseOptions *options.DatabaseOptions) CollectionClient {
	s.WithCollectionReturn(database, collection, databaseOptions)
	return s
}

func (s StubCollectionClient) WithConnection(mongoUri string, opts *options.ClientOptions) CollectionClient {
	s.WithConnectionReturn(mongoUri, opts)
	return s
}

func (s StubCollectionClient) WithLogger(logger *log.Logger) CollectionClient {
	s.WithLoggerReturn(logger)
	return s
}

func (s StubCollectionClient) GetDocumentById(ctx context.Context, recordID interface{}) (MongoCursor, error) {
	return s.GetDocumentByIdReturn(ctx, recordID)
}

func (s StubCollectionClient) SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error) {
	return s.SaveDocumentReturn(ctx, document, modelID)
}

func (s StubCollectionClient) ScanCollection(ctx context.Context) (MongoCursor, error) {
	return s.ScanCollectionReturn(ctx)
}

func (s StubCollectionClient) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return s.CountDocumentsReturn(ctx, filter)
}

func (s StubCollectionClient) CloseConnection(ctx context.Context) error {
	return s.CloseConnectionReturn(ctx)
}

func (s StubCollectionClient) GetMongoClient() *mongo.Client {
	return s.GetMongoClientReturn()
}

func (s StubCollectionClient) GetCollection() *mongo.Collection {
	return s.GetCollectionReturn()
}

func (s StubCollectionClient) MustValidate() {
	s.MustValidateReturn()
}
