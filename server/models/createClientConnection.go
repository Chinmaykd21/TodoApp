package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This function will create a connection with mongoDB atlas and return client, collection.
// NOTE: Since this application will contain only 1 collection I am returning collection in this
// function. Otherwise, this can returned in a separate function.
func CreatConnection(ctx context.Context) (*mongo.Client, *mongo.Collection) {
	// get the mongoDB connection string
	URI := os.Getenv("MONGO_CONNECTION_URL")

	if URI == "" {
		// TODO: should we use panic & recovery instead of log.fatal?
		// TODO: How to return an empty collection object
		log.Fatal("MongoDB connection string cannot be empty")
	}

	// connect to mongodb client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		// TODO: Should we use panic & recovery instead of log.fatal?
		// TODO: How to return an empty collection object
		log.Fatal("Cannot connect to the mongoDB", err)
	}
	fmt.Println("Sucessfully connected to mongoDB atlas")

	// Create a database in mongoDB atlas
	projectsDatabase := client.Database("projects")

	// Create a collection for todo in database projects
	collectionTodos := projectsDatabase.Collection("todos")

	fmt.Println("Successfully created todos collection in projects database")

	return client, collectionTodos
}
