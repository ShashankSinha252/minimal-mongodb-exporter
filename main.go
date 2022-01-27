package main

import (
	"context"
	"flag"
	"fmt"
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

	_, err := connectToMongoDB(l, cfg)
	if err != nil {
		l.Error("Cannot continue ahead: exiting")
		os.Exit(1)
	}

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
