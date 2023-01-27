package mongodb

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

var (
	testId        = primitive.NewObjectID()
	expectedError = errors.New("expected error")
	testModel     = struct {
		ID primitive.ObjectID
	}{
		ID: testId,
	}
)

func TestClient_CountDocuments(t *testing.T) {
	ctx := context.Background()

	t.Run("test when successful", func(t *testing.T) {
		client := MongoClientMock{
			CountDocumentsMockResponse: &CountDocumentsMockResponse{
				Count: 10,
				Error: nil,
			},
			SaveDocumentMockResponse: &SaveDocumentMockResponse{
				Result: true,
				Error:  nil,
			},
			GetDocumentByIdMockResponse: &GetDocumentByIdMockResponse{
				Cursor: &MongoCursorMock{
					DecodeMockResult: DecodeMockResult{
						Error: nil,
					},
				},
			},
		}

		count, err := client.CountDocuments(ctx, primitive.M{"_id": "fake_id"})
		assert.NoError(t, err)

		assert.Equal(t, int64(10), count)
	})

	t.Run("test when error", func(t *testing.T) {
		client := MongoClientMock{
			CountDocumentsMockResponse: &CountDocumentsMockResponse{
				Count: 0,
				Error: expectedError,
			},
			SaveDocumentMockResponse:    nil,
			GetDocumentByIdMockResponse: nil,
		}

		count, err := client.CountDocuments(ctx, primitive.M{"_id": "fake_id"})
		assert.Error(t, err)

		assert.Equal(t, int64(0), count)
	})
}

func TestClient_GetDocumentById(t *testing.T) {
	ctx := context.Background()

	t.Run("test when successful", func(t *testing.T) {
		client := MongoClientMock{
			CountDocumentsMockResponse: nil,
			SaveDocumentMockResponse:   nil,
			GetDocumentByIdMockResponse: &GetDocumentByIdMockResponse{
				Cursor: &MongoCursorMock{
					DecodeMockResult: DecodeMockResult{
						Error: nil,
					},
				},
			},
		}

		_, err := client.GetDocumentById(ctx, testId)
		assert.NoError(t, err)

	})

	t.Run("test when not successful", func(t *testing.T) {
		client := MongoClientMock{
			CountDocumentsMockResponse: nil,
			SaveDocumentMockResponse:   nil,
			GetDocumentByIdMockResponse: &GetDocumentByIdMockResponse{
				Error: expectedError,
				Cursor: &MongoCursorMock{
					DecodeMockResult: DecodeMockResult{
						Error: expectedError,
					},
				},
			},
		}

		_, err := client.GetDocumentById(ctx, testId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedError)

	})

}

func TestMongoClientMock_SaveDocument(t *testing.T) {
	ctx := context.Background()

	t.Run("test when successful", func(t *testing.T) {
		client := MongoClientMock{SaveDocumentMockResponse: &SaveDocumentMockResponse{
			Result: true,
			Error:  nil,
		}}

		isSuccess, err := client.SaveDocument(ctx, testModel, testId)
		assert.NoError(t, err)
		assert.Equal(t, true, isSuccess)
	})

	t.Run("test when not successful", func(t *testing.T) {
		client := MongoClientMock{SaveDocumentMockResponse: &SaveDocumentMockResponse{
			Result: false,
			Error:  expectedError,
		}}

		isSuccess, err := client.SaveDocument(ctx, testModel, testId)
		assert.Error(t, err)
		assert.Equal(t, false, isSuccess)
	})

}
