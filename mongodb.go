package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoDb(host string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}
	return client, nil
}
