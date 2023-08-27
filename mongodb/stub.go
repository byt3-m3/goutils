package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	SaveDocumentStubResponse struct {
		Result bool
		Error  error
	}
	CountDocumentsStubResponse struct {
		Count int64
		Error error
	}

	GetDocumentByIdStubResponse struct {
		Cursor MongoCursor
		Error  error
	}

	CloseConnectionStubResponse struct {
		Error error
	}
)

type StubMongoClient struct {
	CountDocumentsStubResponse   func(ctx context.Context, filter primitive.M) CountDocumentsStubResponse
	SaveDocumentStubResponse     func(ctx context.Context, document interface{}, modelID interface{}) SaveDocumentStubResponse
	GetDocumentByIdStubResponse  func(ctx context.Context, recordId interface{}) GetDocumentByIdStubResponse
	CloseConnectionStubResponse  func(ctx context.Context) CloseConnectionStubResponse
	GetMongoCollectionStubReturn func() GetMongoCollectionStubReturn
	GetMongoClientStubReturn     func() GetMongoClientStubReturn
	ScanCollectionStubReturn     func(ctx context.Context) ScanCollectionStubReturn
}

type ScanCollectionStubReturn struct {
	Cursor MongoCursor
	Err    error
}

func (m *StubMongoClient) ScanCollection(ctx context.Context) (MongoCursor, error) {
	res := m.ScanCollectionStubReturn(ctx)
	return res.Cursor, res.Err
}

type GetMongoCollectionStubReturn struct {
	Collection *mongo.Collection
}

func (m *StubMongoClient) GetMongoCollection() *mongo.Collection {
	return m.GetMongoCollectionStubReturn().Collection
}

type GetMongoClientStubReturn struct {
	Client *mongo.Client
}

func (m *StubMongoClient) GetMongoClient() *mongo.Client {
	return m.GetMongoClientStubReturn().Client
}

func (m *StubMongoClient) GetDocumentById(ctx context.Context, recordId interface{}) (MongoCursor, error) {
	res := m.GetDocumentByIdStubResponse(ctx, recordId)
	return res.Cursor, res.Error
}

func (m *StubMongoClient) SaveDocument(ctx context.Context, document interface{}, modelID interface{}) (bool, error) {
	res := m.SaveDocumentStubResponse(ctx, document, modelID)
	return res.Result, res.Error
}

func (m *StubMongoClient) CountDocuments(ctx context.Context, filter primitive.M) (int64, error) {
	res := m.CountDocumentsStubResponse(ctx, filter)
	return res.Count, res.Error
}
func (m *StubMongoClient) CloseConnection(ctx context.Context) error {
	return m.CloseConnectionStubResponse(ctx).Error
}
