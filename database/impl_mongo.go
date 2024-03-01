package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	collection *mongo.Collection
}

type URL struct {
	Long   string `json:"long"`
	Short  string `json:"short"`
	Clicks int    `json:"clicks"`
}

type Database interface {
	InsertOne(url URL) error
	FindOne(shorturl string) (*URL, error)
	Find() ([]URL, error)
	Update(url URL) (*URL, error)
}

func InitMongo(databaseURL string) (Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(databaseURL))
	if err != nil {
		return nil, err
	}

	collection := client.Database("urlshortener").Collection("urls")

	return &mongoDB{collection}, nil
}

func (e *mongoDB) InsertOne(url URL) error {
	_, err := e.collection.InsertOne(context.Background(), url)
	if err != nil {
		return err
	}

	return nil
}

func (e *mongoDB) FindOne(shorturl string) (*URL, error) {

	var url URL
	err := e.collection.FindOne(context.TODO(), bson.M{"short": shorturl}).Decode(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (e *mongoDB) Find() ([]URL, error) {

	var urls []URL

	cursor, err := e.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var url URL
		if err := cursor.Decode(&url); err != nil {
			continue
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (e *mongoDB) Update(url URL) (*URL, error) {

	update := bson.M{
		"$set": bson.M{
			"long":   url.Long,
			"short":  url.Short,
			"clicks": url.Clicks,
		},
	}

	_, err := e.collection.UpdateOne(context.TODO(), bson.M{"short": url.Short}, update)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
