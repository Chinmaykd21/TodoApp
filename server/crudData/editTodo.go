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

	todoToEdit := todo
	todoAfterUpdate := []customDataStructs.Todo{}

	// Filter
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// Update condition
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: todoToEdit.Title},
		{Key: "body", Value: todoToEdit.Body}}}}

	// Update query
	updateResult, err := collectionTodos.UpdateOne(ctx, filter, update)

	if err != nil {
		errResponse := serverErrors.New(serverErrors.UpdateError, err.Error())
		return &todoAfterUpdate, errResponse
	}

	if updateResult.ModifiedCount == 0 {
		errResponse := serverErrors.New(serverErrors.UpdateError, "Modified records are 0. No Records have been updated")
		return &todoAfterUpdate, errResponse
	}

	if updateResult.MatchedCount != updateResult.ModifiedCount {
		errResponse := serverErrors.New(serverErrors.UpdateError, "Matched & modified records have different counts")
		return &todoAfterUpdate, errResponse
	}

	todoAfterUpdate = append(todoAfterUpdate, todoToEdit)

	return &todoAfterUpdate, nil
}
