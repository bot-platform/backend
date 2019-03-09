package main

import (
	"context"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
}

func (repo *Repository) InsertTeam(team *Team) error {
	res, err := repo.db.Collection("teams").InsertOne(context.TODO(), team)
	if err != nil {
		return err
	}
	team.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (repo *Repository) FindTournaments() ([]*Tournament, error) {
	var results []*Tournament
	cur, err := repo.db.Collection("tournaments").Find(context.TODO(), nil, options.Find())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	for cur.Next(context.TODO()) {
		var elem Tournament
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *Repository) FindBots() ([]*Bot, error) {
	var results []*Bot
	cur, err := repo.db.Collection("bots").Find(context.TODO(), nil, options.Find())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cur.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	for cur.Next(context.TODO()) {
		var elem Bot
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{
		client: client,
		db:     client.Database(dbName),
	}
}
