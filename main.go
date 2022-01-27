package main

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {

	// Initial setup
	l := hclog.New(nil)

	_, err := connectToMongoDB(l)
	if err != nil {
		l.Error("Cannot continue ahead: exiting")
		os.Exit(1)
	}

}

func connectToMongoDB(l hclog.Logger) (*mongo.Client, error) {
	dsn := "mongodb://localhost:27017"

	// Set client options
	clientOptions := options.Client().ApplyURI(dsn)

	// Connect to MongoDB
	l.Info("Attempting to connect to MongoDB", "dsn", dsn)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		l.Error("Failed to connect", "error msg", err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		l.Error("Ping failed for MongoDB client", "error msg", err)
		return nil, err
	}

	l.Info("Connected to MongoDB!")

	return client, nil
}
