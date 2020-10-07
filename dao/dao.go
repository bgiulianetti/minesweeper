package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mercadolibre/minesweeper/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBContainer ...
type MongoDBContainer struct {
	Client     *mongo.Client
	DB         string
	Collection string
}

// CreateContainer initialize the container
func CreateContainer() MongoDBContainer {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://bgiulianetti:mongodb@cluster0.0nvl0.mongodb.net/minesweper?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	mongoDBContainer := MongoDBContainer{
		Client:     client,
		DB:         "minesweper",
		Collection: "games",
	}
	return mongoDBContainer
}

// GetAll gets all games
func (mdb *MongoDBContainer) GetAll() ([]*domain.UserGame, error) {

	var userGames []*domain.UserGame
	collection := mdb.Client.Database(mdb.DB).Collection(mdb.Collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var userGame domain.UserGame
		cursor.Decode(&userGame)
		userGames = append(userGames, &userGame)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return userGames, nil
}

// Get a userGame by userID
func (mdb *MongoDBContainer) Get(userID string) (*domain.UserGame, error) {

	var userGame *domain.UserGame
	collection := mdb.Client.Database(mdb.DB).Collection(mdb.Collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userGame)
	if err != nil {
		return nil, err
	}
	return userGame, nil
}

// Update ...
func (mdb *MongoDBContainer) Update(userGame *domain.UserGame) error {

	collection := mdb.Client.Database(mdb.DB).Collection(mdb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.ReplaceOne(
		ctx,
		bson.M{"user_id": userGame.UserID},
		bson.M{
			"user_id": userGame.UserID,
			"games":   userGame.Games,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Insert ...
func (mdb *MongoDBContainer) Insert(userGame *domain.UserGame) error {
	collection := mdb.Client.Database(mdb.DB).Collection(mdb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, userGame)
	if err != nil {
		return err
	}
	fmt.Println(res.InsertedID)
	return nil
}

// DeleteAll deletes all games
func (mdb *MongoDBContainer) DeleteAll() error {

	collection := mdb.Client.Database(mdb.DB).Collection(mdb.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
