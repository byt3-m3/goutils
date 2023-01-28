package mongodb

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/dig"
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

	testClientConfig = &ClientConfig{
		MClientConfig:  &MongoClientConfig{MongoURI: "mongodb://192.168.1.5"},
		CollectionName: "test-collection",
		DBName:         "test-db",
	}

	mockClientSuccessProvider = func() IMongoClient {
		return &MongoClientMock{SaveDocumentMockResponse: &SaveDocumentMockResponse{
			Result: true,
			Error:  nil,
		}}

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

	c := dig.New()

	if err := c.Provide(mockClientSuccessProvider); err != nil {
		t.Fatal(err)
	}

	//// uncomment to do live test
	//if err := c.Provide(func() *ClientConfig {
	//	return testClientConfig
	//}); err != nil {
	//	t.Fatal(err)
	//}
	//
	//if err := c.Provide(ProvideIMongoClient); err != nil {
	//	t.Fatal(err)
	//}

	if err := c.Invoke(func(client IMongoClient) {

		t.Run("test when successful", func(t *testing.T) {

			isSuccess, err := client.SaveDocument(ctx, testModel, testId)
			assert.NoError(t, err)
			assert.Equal(t, true, isSuccess)
		})

	}); err != nil {
		t.Fatal(err)
	}

}
