package crudData

import (
	"context"
	"net/http"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Function to add the new todo
func AddToDo(ctx context.Context, c *fiber.Ctx, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) error {
	// create an empty todo
	todo := &customDataStructs.Todo{}

	// If it returns an error, return that error
	if err := c.BodyParser(todo); err != nil {
		errResponse := serverErrors.New(serverErrors.BodyParse, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return err
	}

	// new unique id for the todo
	todoId := len(todos) + 1

	// insert first dummy todo document to the collection
	_, err := collectionTodos.InsertOne(ctx, bson.D{
		{Key: "todoId", Value: todoId},
		{Key: "title", Value: &todo.Title},
		{Key: "body", Value: &todo.Body},
		{Key: "isCompleted", Value: &todo.IsCompleted},
	})

	if err != nil {
		// TODO: should we use panic & recovery instead of log.fatal?
		errResponse := serverErrors.New(serverErrors.InsertError, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return err
	}

	return nil
}
