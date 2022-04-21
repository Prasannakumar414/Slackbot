package main

import (
	"context"

	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func UpsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}, name string) (*mongo.UpdateResult, error) {

	collection := client.Database(dataBase).Collection(col)

	filter := bson.D{{"Name", name}}
	opts := options.Replace().SetUpsert(true)
	result, err := collection.ReplaceOne(context.TODO(), filter, doc, opts)
	return result, err
}

func query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := client.Database(dataBase).Collection(col)

	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return result, err
}

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
