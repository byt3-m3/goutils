package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"strings"
)

type (
	MongoClientConfig struct {
		MongoURI string
	}

	SaveOrUpdateModelResult struct {
		IsSuccess bool
		ModelId   interface{}
	}
)

type FindModelResult struct {
	MongoCursor *mongo.Cursor
	HasData     bool
}

func FindDocument(ctx context.Context, docFinder MongoCollectionFinder, filter primitive.M) (FindModelResult, error) {
	cur, err := docFinder.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return FindModelResult{HasData: false, MongoCursor: cur}, DocumentNotFoundError
		}

		return FindModelResult{HasData: false}, err
	}
	return FindModelResult{HasData: true, MongoCursor: cur}, nil
}

func SaveOrUpdateDocument(ctx context.Context, collection *mongo.Collection, model interface{}, modelId interface{}) (*SaveOrUpdateModelResult, error) {
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		if strings.Contains(err.Error(), "dup key") {

			if replaceErr := ReplaceDocumentById(ctx, collection, model, modelId); err != nil {
				return &SaveOrUpdateModelResult{IsSuccess: true, ModelId: modelId}, replaceErr
			}

		}
		return nil, err
	}
	if result.InsertedID != nil {

		return &SaveOrUpdateModelResult{IsSuccess: true, ModelId: result.InsertedID}, nil
	}
	return nil, nil

}

func ReplaceDocument(ctx context.Context, collection *mongo.Collection, model interface{}, filter primitive.M) error {
	updateResult, replaceErr := collection.ReplaceOne(
		ctx,
		filter,
		model,
	)

	if replaceErr != nil {
		return replaceErr
	}
	if updateResult.ModifiedCount > 0 {
		slog.Debug("replace successful")

	} else {
		slog.Debug("no records modified")
	}

	return nil

}

func ReplaceDocumentById(ctx context.Context, collection *mongo.Collection, model interface{}, modelId interface{}) error {
	return ReplaceDocument(ctx, collection, model, primitive.M{"_id": modelId})

}

func DeleteDocument(ctx context.Context, cClient CollectionClient, filter interface{}, logger slog.Logger) error {
	return cClient.DeleteDocument(ctx, filter, logger)
}
