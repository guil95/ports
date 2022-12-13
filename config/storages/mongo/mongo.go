package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func Connect(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", os.Getenv("MONGODB_HOST")))
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		zap.S().Fatal(err)
		return nil
	}

	return client
}
