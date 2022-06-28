package crudData

import (
	"context"
	"fmt"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllTodos(ctx context.Context, c *fiber.Ctx, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	cursor, err := collectionTodos.Find(ctx, bson.M{})
	if err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, err.Error())
		return &todos, errResponse
	}

	defer func() {
		fmt.Println("Closing cursor")
		cursor.Close(ctx)
	}()

	if err = cursor.All(ctx, &todos); err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, err.Error())
		return &todos, errResponse
	}

	return &todos, nil
}
