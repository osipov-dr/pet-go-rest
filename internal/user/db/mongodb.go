package db

import (
	"context"
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
		//TODO 404
		return user, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode user (id: %s) from db due to error: %v", id, err)
	}

	return user, nil
}

func (db *db) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (db *db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
