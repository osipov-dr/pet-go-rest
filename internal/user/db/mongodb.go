package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"osipovPetRestApi/internal/user"
	"osipovPetRestApi/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (db *db) Create(ctx context.Context, user user.User) (userId string, err error) {
	db.logger.Debug("create user")
	result, err := db.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("new user creation failed: %v", err)
	}

	db.logger.Debug("convert InsertedId to ObjectId")
	objectId, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return objectId.Hex(), nil
	}
	db.logger.Trace(user)

	return "", fmt.Errorf("objectId to hex convertion failed, objectId: %s", objectId)
}

func (db *db) FindOne(ctx context.Context, id string) (user user.User, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, fmt.Errorf("hex to objectId convertion failed, hex: %s", id)
	}

	filter := bson.M{"_id": objectId}
	result := db.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
			return user, fmt.Errorf("not found")
		}
		return user, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode user (id: %s) from db due to error: %v", id, err)
	}

	return user, nil
}

func (db *db) Update(ctx context.Context, user user.User) error {
	objectId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return fmt.Errorf("hex to objectId convertion failed, hex: %s", user.Id)
	}

	filter := bson.M{"_id": objectId}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObject bson.M
	err = bson.Unmarshal(userBytes, &updateUserObject)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. error: %v", err)
	}

	delete(updateUserObject, "_id")

	update := bson.M{
		"$set": updateUserObject,
	}
	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %w", err)
	}

	if result.MatchedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	db.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (db *db) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("hex to objectId convertion failed, hex: %s", id)
	}

	filter := bson.M{"_id": objectId}

	result, err := db.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}

	if result.DeletedCount == 0 {
		// TODO ErrEntityNotFound
		return fmt.Errorf("not found")
	}
	db.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
