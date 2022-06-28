package crudData

import (
	"context"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditTodo(ctx context.Context, c *fiber.Ctx, todoId int, todo customDataStructs.Todo, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	var todoToEdit customDataStructs.Todo
	var editedTodoInSlice []customDataStructs.Todo

	// Filter
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// Find and store the id whose status is to be toggled to struct todo
	err := collectionTodos.FindOne(ctx, filter).Decode(&todoToEdit)
	if err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, err.Error())
		return &editedTodoInSlice, errResponse
	}

	// Update condition
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: todoToEdit.Title},
		{Key: "body", Value: todoToEdit.Body},
		{Key: "isCompleted", Value: todoToEdit.IsCompleted}}}}

	// Update query
	updateResult, err := collectionTodos.UpdateOne(ctx, filter, update)

	if err != nil {
		errResponse := serverErrors.New(serverErrors.UpdateError, err.Error())
		return &editedTodoInSlice, errResponse
	}

	if updateResult.MatchedCount != updateResult.ModifiedCount {
		errResponse := serverErrors.New(serverErrors.UpdateError, err.Error())
		return &editedTodoInSlice, errResponse
	}

	editedTodoInSlice = append(editedTodoInSlice, todoToEdit)

	return &editedTodoInSlice, nil
}
