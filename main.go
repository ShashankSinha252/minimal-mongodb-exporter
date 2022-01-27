package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

type Config struct {
	host     string
	port     int
	username string
	password string
}

const (
	defaultHost = "localhost"
	defaultPort = 27017
)

func parseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.host, "net.host", defaultHost, "host running MongoDB")
	flag.IntVar(&cfg.port, "net.port", defaultPort, "port for MongoDB")
	flag.StringVar(&cfg.username, "db.username", "", "user for MongoDB")
	flag.StringVar(&cfg.password, "db.password", "", "password for MongoDB")

	flag.Parse()
	return cfg
}

func main() {

	// Initial setup
	l := hclog.New(nil)
	cfg := parseFlags()

	client, err := connectToMongoDB(l, cfg)
	if err != nil {
		l.Error("Cannot continue ahead: exiting")
		os.Exit(1)
	}

	collection := client.Database("test").Collection("trainers")

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	insertInCollection(l, collection, ash)
	insertInCollection(l, collection, misty)
	insertInCollection(l, collection, brock)
	updateInCollection(l, collection)
	findInCollection(l, collection)

	disconnectFromMongoDB(l, client)
}

// createDSN generates a DSN for mongoDB
func createDSN(host, username, password string, port int) string {
	var dsn string
	if username == "" {
		dsn = fmt.Sprintf("mongodb://%s:%d", host, port)
	} else {
		dsn = fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	}

	return dsn
}

// connectToMongoDB uses given params to create a MongoDB client
func connectToMongoDB(l hclog.Logger, params *Config) (*mongo.Client, error) {
	dsn := createDSN(params.host, params.username, params.password, params.port)

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

// disconnectFromMongoDB will disconnect given MongoDB client
func disconnectFromMongoDB(l hclog.Logger, client *mongo.Client) {

	err := client.Disconnect(context.TODO())
	if err != nil {
		l.Error("Error while disconnecting from MongoDB", "error", err)
	}

	l.Info("Connection to MongoDB closed")
}

// insertInCollection inserts given document in given collection
func insertInCollection(l hclog.Logger, col *mongo.Collection, doc interface{}) error {

	insertResult, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		l.Error("Unable to insert document", "error", err)
		return err
	}

	l.Info("Inserted a single document: ", "ID", insertResult.InsertedID)
	return nil

}

func updateInCollection(l hclog.Logger, col *mongo.Collection) error {
	filter := bson.M{"name": "Ash"}

	update := bson.M{
		"$inc": bson.M{
			"age": 1,
		},
	}

	updateResult, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		l.Error("Unable to update document", "error", err)
		return err
	}

	l.Info("Updated documents", "Matched documents", updateResult.MatchedCount, "Updated documents", updateResult.ModifiedCount)
	return nil
}

func findInCollection(l hclog.Logger, col *mongo.Collection) error {
	var result Trainer
	filter := bson.M{"name": "Ash"}

	err := col.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		l.Error("Unable to find document", "error", err)
		return err
	}

	l.Info("Found document", "Document", result)
	return nil
}
