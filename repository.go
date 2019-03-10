package main

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
}

func (repo *Repository) InsertTeam(team *Team) error {
	team.ID = primitive.NewObjectID()
	_, err := repo.db.Collection("teams").InsertOne(context.TODO(), team)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) InsertBot(bot *Bot) error {
	bot.ID = primitive.NewObjectID()
	_, err := repo.db.Collection("bots").InsertOne(context.TODO(), bot)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) FindTournaments() ([]*Tournament, error) {
	var results = make([]*Tournament, 0)
	cur, err := repo.db.Collection("tournaments").Find(context.TODO(), bson.D{}, options.Find())
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
	find := options.Find()
	find.Sort = bson.M{
		"_id": -1,
	}
	var results = make([]*Bot, 0)
	cur, err := repo.db.Collection("bots").Find(context.TODO(), bson.D{}, find)
	if err != nil {
		return nil, errors.Wrap(err, "find query failed")
	}
	defer func() {
		if err := cur.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	for cur.Next(context.TODO()) {
		var bot Bot
		err := cur.Decode(&bot)
		if err != nil {
			return nil, errors.Wrap(err, "cursor decode error")
		}
		results = append(results, &bot)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor error")
	}
	return results, nil
}

func (repo *Repository) GetUserByLogin(login string) (*User, error) {
	var user User
	if err := repo.db.Collection("users").FindOne(context.TODO(), bson.M{"login": login}, options.FindOne()).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "decode find one query failed")
	}
	return &user, nil
}

func (repo *Repository) InsertUser(user *User) error {
	user.ID = primitive.NewObjectID()
	_, err := repo.db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetBotByID(id string) (*Bot, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "cant transform string to ObjectID")
	}
	var bot Bot
	if err := repo.db.Collection("bots").FindOne(context.TODO(), bson.M{"_id": objectID}, options.FindOne()).Decode(&bot); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "decode find one query failed")
	}
	return &bot, nil
}

func (repo *Repository) DeleteBot(bot *Bot) error {
	res, err := repo.db.Collection("bots").DeleteOne(context.TODO(), bson.M{"_id": bot.ID}, options.Delete())
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("zero documents were deleted")
	}
	return nil
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{
		client: client,
		db:     client.Database(dbName),
	}
}
