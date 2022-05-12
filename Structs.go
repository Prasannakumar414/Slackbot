package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	client *mongo.Client
	cancel context.CancelFunc
}
type Standup struct {
	ID               string   `bson:"_id" json:"ID"`
	Name             string   `bson:"Name" json:"Name"`
	Hour             int      `bson:"Hour" json:"Hour"`
	Minute           int      `bson:"Minute" json:"Minute"`
	Days             []int    `bson:"Days" json:"Days"`
	BroadcastChannel string   `bson:"BroadcastChannel" json:"BroadcastChannel"`
	Users            []string `bson:"Users" json:"Users"`
}

type Response struct {
	Sender  string `bson:"Sender" json:"Sender"`
	Message string `bson:"Message" json:"Message"`
}

func NewServer(client *mongo.Client, cancel context.CancelFunc) *Server {
	return &Server{
		client: client,
		cancel: cancel,
	}
}
