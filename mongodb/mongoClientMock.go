package mongodb

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	SaveDocumentMockResponse struct {
		Result bool
		Error  error
	}
	CountDocumentsMockResponse struct {
		Count int64
		Error error
	}

	GetDocumentByIdMockResponse struct {
		Cursor MongoCursor
		Error  error
	}
)

type MongoClientMock struct {
	mock.Mock
	CountDocumentsMockResponse  *CountDocumentsMockResponse
	SaveDocumentMockResponse    *SaveDocumentMockResponse
	GetDocumentByIdMockResponse *GetDocumentByIdMockResponse
}

func (m *MongoClientMock) GetDocumentById(ctx context.Context, recordId primitive.ObjectID) (MongoCursor, error) {
	return m.GetDocumentByIdMockResponse.Cursor, m.GetDocumentByIdMockResponse.Error
}

func (m *MongoClientMock) SaveDocument(ctx context.Context, document interface{}, modelID primitive.ObjectID) (bool, error) {
	return m.SaveDocumentMockResponse.Result, m.SaveDocumentMockResponse.Error
}

func (m *MongoClientMock) CountDocuments(ctx context.Context, filter primitive.M) (int64, error) {
	return m.CountDocumentsMockResponse.Count, m.CountDocumentsMockResponse.Error
}
