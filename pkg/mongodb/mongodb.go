package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/shalimski/shortener/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient established connection to a MongoDB instance using provided URI and auth credentials.
func NewClient(cfg *config.Config) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)
	opts := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      cfg.Mongo.User,
		Password:      cfg.Mongo.Password,
	})

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return client, nil
}
