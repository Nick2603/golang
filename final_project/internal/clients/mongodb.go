package clients

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	slog.Info("Connected to MongoDB", "database", dbName)

	return &MongoDB{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.db.Collection(name)
}

func (m *MongoDB) Disconnect(ctx context.Context) {
	if err := m.client.Disconnect(ctx); err != nil {
		slog.Error("Failed to disconnect from MongoDB", "error", err)
	} else {
		slog.Info("Disconnected from MongoDB")
	}
}
