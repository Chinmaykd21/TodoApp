package crudData

import (
	"context"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllTodos(ctx context.Context, c *fiber.Ctx, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	cursor, err := collectionTodos.Find(ctx, bson.M{})
	if err != nil {
		return &todos, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &todos); err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, "")
		// c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return &todos, err
	}

	return &todos, nil
}
