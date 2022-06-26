package crudData

import (
	"context"

	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// To delete a task based on todoId
func DeleteTodo(ctx context.Context, c *fiber.Ctx, todoId int, collectionTodos *mongo.Collection) error {
	// Creating the filter to get the document which we want to delete from todos collection
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// To query collectionTodos and delete the document
	deletedResult, err := collectionTodos.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if deletedResult.DeletedCount == 0 {
		errResponse := serverErrors.New(serverErrors.RecordsNotFound, "")
		return errResponse
	}

	return nil
}
