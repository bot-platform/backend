package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tournament struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Image       string             `json:"image" bson:"image"`
	Description string             `json:"description" bson:"description"`
}

type Bot struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Address string             `json:"address" bson:"address"`
}

type Team struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Login    string             `json:"login"`
	Password string             `json:"-"`
}

type Player struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type Game struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
}

type Goal struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Key         string             `json:"key"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Score       int                `json:"score"`
}
