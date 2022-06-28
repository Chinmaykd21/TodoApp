package crudData

import (
	"context"

	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteTodo(ctx context.Context, c *fiber.Ctx, todoId int, collectionTodos *mongo.Collection) error {
	// Filter to select only documents matching with todoId
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// Delete the first matching todo from the todos mongoDB collection
	deletedResult, err := collectionTodos.DeleteOne(ctx, filter)

	if err != nil {
		errResponse := serverErrors.New(serverErrors.RecordsNotFound, err.Error())
		return errResponse
	}

	if deletedResult.DeletedCount == 0 {
		errResponse := serverErrors.New(serverErrors.RecordsNotFound, err.Error())
		return errResponse
	}

	return nil
}
