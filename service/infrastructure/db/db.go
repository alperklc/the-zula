package db

import (
	"context"
	"time"

	"github.com/alperklc/the-zula/service/infrastructure/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(mongoUri string) *mongo.Database {
	duration := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	l := logger.Get()
	clientOptions := options.Client().ApplyURI(mongoUri).SetDirect(true)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.Fatal().Msgf("could not connect to MongoDB, %s", err)
	}

	l.Info().Msgf("Connected to MongoDB %s", mongoUri)
	db := mongoClient.Database("zula")
	return db
}
