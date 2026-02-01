package db

import (
	"context"
	"fmt"
	"time"

	"todo-api-go/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func ConnectMongo(cfg config.Config) (*mongo.Client, *mongo.Collection, error) {
	if cfg.MongoURI == "" {
		return nil, nil, fmt.Errorf("MONGO_URI is required")
	}

	clientOpts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetMonitor(otelmongo.NewMonitor()) // <-- OTEL for Mongo

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, nil, err
	}

	// Ping
	if err := client.Database("admin").RunCommand(ctx, map[string]any{"ping": 1}).Err(); err != nil {
		return nil, nil, err
	}

	coll := client.Database(cfg.MongoDB).Collection(cfg.MongoCollection)
	return client, coll, nil
}
