package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mockClientSuccessProvider = func() MongoClient {
		return &StubMongoClient{
			SaveDocumentStubResponse: func(ctx context.Context, document interface{}, modelID interface{}) SaveDocumentStubResponse {
				return SaveDocumentStubResponse{
					Result: true,
					Error:  nil,
				}
			},
		}
	}
)

func TestClient_CountDocuments(t *testing.T) {
	ctx := context.Background()

	t.Run("test when successful", func(t *testing.T) {
		client := StubMongoClient{
			CountDocumentsStubResponse: func(ctx context.Context, filter primitive.M) CountDocumentsStubResponse {

				return CountDocumentsStubResponse{
					Count: 10,
					Error: nil,
				}
			},
			SaveDocumentStubResponse: func(ctx context.Context, document interface{}, modelID interface{}) SaveDocumentStubResponse {

				return SaveDocumentStubResponse{
					Result: true,
					Error:  nil,
				}
			},
			GetDocumentByIdStubResponse: func(ctx context.Context, recordId interface{}) GetDocumentByIdStubResponse {

				return GetDocumentByIdStubResponse{
					Cursor: &StubMongoCursor{
						DecodeStubResult: func(val interface{}) DecodeStubResult {
							return DecodeStubResult{
								Error: nil,
							}
						},
					},
				}
			},
		}

		count, err := client.CountDocuments(ctx, primitive.M{"_id": "fake_id"})
		assert.NoError(t, err)

		assert.Equal(t, int64(10), count)
	})

	t.Run("test when error", func(t *testing.T) {
		client := StubMongoClient{
			CountDocumentsStubResponse: func(ctx context.Context, filter primitive.M) CountDocumentsStubResponse {
				return CountDocumentsStubResponse{
					Count: 0,
					Error: expectedError,
				}
			},
			SaveDocumentStubResponse:    nil,
			GetDocumentByIdStubResponse: nil,
		}

		count, err := client.CountDocuments(ctx, primitive.M{"_id": "fake_id"})
		assert.Error(t, err)

		assert.Equal(t, int64(0), count)
	})
}

func TestClient_GetDocumentById(t *testing.T) {
	ctx := context.Background()

	t.Run("test when successful", func(t *testing.T) {
		client := StubMongoClient{
			CountDocumentsStubResponse: nil,
			SaveDocumentStubResponse:   nil,
			GetDocumentByIdStubResponse: func(ctx context.Context, recordId interface{}) GetDocumentByIdStubResponse {
				return GetDocumentByIdStubResponse{
					Cursor: &StubMongoCursor{
						DecodeStubResult: func(val interface{}) DecodeStubResult {
							return DecodeStubResult{}
						},
					},
				}
			},
		}

		_, err := client.GetDocumentById(ctx, testId)
		assert.NoError(t, err)

	})

	t.Run("test when not successful", func(t *testing.T) {
		client := StubMongoClient{
			CountDocumentsStubResponse: nil,
			SaveDocumentStubResponse:   nil,
			GetDocumentByIdStubResponse: func(ctx context.Context, recordId interface{}) GetDocumentByIdStubResponse {
				return GetDocumentByIdStubResponse{
					Error: expectedError,
					Cursor: &StubMongoCursor{
						DecodeStubResult: func(val interface{}) DecodeStubResult {
							return DecodeStubResult{
								Error: expectedError,
							}
						},
					},
				}
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

	if err := c.Invoke(func(client MongoClient) {

		t.Run("test when successful", func(t *testing.T) {

			isSuccess, err := client.SaveDocument(ctx, testModel, testId)
			assert.NoError(t, err)
			assert.Equal(t, true, isSuccess)
		})

	}); err != nil {
		t.Fatal(err)
	}

}

func TestName(t *testing.T) {
	c := NewClient(
		WithMongoClient("mongodb://root:mongo@192.168.1.58", &options.ClientOptions{}),
		WithCollection("test-db", "test-collection"),
	)
	count, err := c.SaveDocument(context.TODO(), map[string]string{"test": "test"}, "1")

	fmt.Println(err, count)
}
