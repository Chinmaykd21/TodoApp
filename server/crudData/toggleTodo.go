package crudData

import (
	"context"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateTodo(ctx context.Context, c *fiber.Ctx, todoId int, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	// Initialize an empty todo
	var todoToUpdate customDataStructs.Todo
	var updatedTodoInSlice []customDataStructs.Todo

	// Define filter condition to find the document with proper val
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// Find and store the id whose status is to be toggled to struct todo
	err := collectionTodos.FindOne(ctx, filter).Decode(&todoToUpdate)
	if err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, "")
		return &updatedTodoInSlice, errResponse
	}

	// Toggle its current state & store it
	toggleIsCompleted := !todoToUpdate.IsCompleted

	// Update condition
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isCompleted", Value: toggleIsCompleted}}}}

	// Update the same ID with the new state
	updateResult, err := collectionTodos.UpdateOne(ctx, filter, update)

	if err != nil {
		errResponse := serverErrors.New(serverErrors.UpdateError, "")
		return &updatedTodoInSlice, errResponse
	}

	if updateResult.MatchedCount != updateResult.ModifiedCount {
		errResponse := serverErrors.New(serverErrors.UpdateError, "")
		return &updatedTodoInSlice, errResponse
	}

	updatedTodoInSlice = append(updatedTodoInSlice, todoToUpdate)

	return &updatedTodoInSlice, nil
}
